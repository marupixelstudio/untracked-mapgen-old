mapnik = require('mapnik')
fs = require('fs')
mapnik.register_default_fonts()
mapnik.register_default_input_plugins()
map = new (mapnik.Map)(1200, 800)

map.load 'styles/layers-contours.xml', (err, map) ->
  throw err if err
  map.zoomAll()
  img = new (mapnik.Image)(1200, 800)
  map.render img, (err, img) ->
    throw err if err
    img.encode 'png', (err, buffer) ->
      throw err if err
      fs.writeFile 'map.png', buffer, (err) ->
        throw err if err
        console.log 'saved map image to map.png'
