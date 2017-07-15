# BrickRecon
Helps find Lego Bricks



## Stories

> Given a BSX file of a model, Display a web page of all bricks within, with: names, partno, colour, quantity, and image



## Initial Design

* S3 bucket
  * with folders: `upload`, `website\models`, `website\images\parts`
* Static website
  * bootstrap webpage
  * table contents from url querysting (e.g. `?model=fulcrum`)
  * call to `website\models\fulcrum.json`
  * render into a table in page
* Lambda: on upload
  * triggered by upload/change within `upload` folder
  * parses bsx file
  * for each part
    * build row for table
      * generate url for image, it might not exist yet
    * add to image lookup batch
  * trigger image lookup
  * render json to `website\models` folder, using bsx filename
* Lambda: image lookup
  * bricklink images look to be addressable
    * default: https://img.bricklink.com/ItemImage/PN/1/3022.png
    * dark stone grey: https://img.bricklink.com/ItemImage/PN/85/3022.png
    * dark stone grey: https://img.bricklink.com/ItemImage/PN/{colorCode}/{partno}.png
    * size is consistent for colorCode > 1, but not across parts
    * might want to autocrop images on receive
  * write images to `website\images\parts`
