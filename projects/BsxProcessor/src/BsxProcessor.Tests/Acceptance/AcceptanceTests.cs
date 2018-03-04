using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.Model;
using Amazon.Lambda.SNSEvents;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;
using Newtonsoft.Json;
using NSubstitute;
using Shouldly;
using Xunit;

namespace BsxProcessor.Tests.Acceptance
{
	public class AcceptanceTests
	{
		private SnsHandler _snsHandler;
		private InMemoryFileSystem _fileSystem;

		public AcceptanceTests()
		{
			var config = new Config
			{
				ImageCacheLambda = "wat"
			};

			var imageRequests = new List<InvokeRequest>();

			var imageCacheDispatch = new ImageCacheDispatcher(config, req =>
			{
				imageRequests.Add(req);
				return Task.FromResult(new InvokeResponse());
			});
			var modelBuilder = new BsxModelBuilder();

			_fileSystem = new InMemoryFileSystem();
			var bsxProcessor = new BsxProcessor(_fileSystem, imageCacheDispatch, modelBuilder);

			_snsHandler = new SnsHandler(bsxProcessor);
		}

		[Fact]
		public async Task When_called_from_sns()
		{
			var sns = JsonConvert.DeserializeObject<SNSEvent>(await ReadJson());

			await _snsHandler.Handle(sns);

			var written = _fileSystem.Writes.OfType<FileData<BsxModel>>().Single();

			written.ShouldSatisfyAllConditions(
				() => written.Drive.ShouldBe("brickrecon-dev"),
				() => written.FullPath.ShouldBe("models/testmodel.json")
			);
		}

		private async Task<string> ReadJson()
		{
			var assembly = GetType().GetTypeInfo().Assembly;

			using (var ms = assembly.GetManifestResourceStream("BsxProcessor.Tests.Resources.SnsEvent.json"))
			using (var reader = new StreamReader(ms))
			{
				return await reader.ReadToEndAsync();
			}
		}

		public class InMemoryFileSystem : IFileSystem
		{
			public IEnumerable<object> Writes => _writes;

			private readonly List<object> _writes;

			public InMemoryFileSystem()
			{
				_writes = new List<object>();
			}

			public Task<FileData<XDocument>> ReadXml(string drive, string path)
			{
				throw new System.NotImplementedException();
			}

			public Task WriteJson<TContent>(FileData<TContent> file)
			{
				_writes.Add(file);
				return Task.CompletedTask;
			}
		}
	}
}
