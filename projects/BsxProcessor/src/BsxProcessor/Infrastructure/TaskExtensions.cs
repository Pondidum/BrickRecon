using System;
using System.Threading.Tasks;

namespace BsxProcessor.Infrastructure
{
	public static class TaskExtensions
	{
		public static async Task<TOut> Start<TIn, TOut>(this TIn value, Func<TIn, Task<TOut>> next)
		{
			return await next(value);
		}

		public static async Task<TOut> Then<TIn, TOut>(this Task<TIn> task, Func<TIn, Task<TOut>> next)
		{
			return await next(await task);
		}

		public static async Task Then<TIn>(this Task<TIn> task, Func<TIn, Task> next)
		{
			await next(await task);
		}
	}
}