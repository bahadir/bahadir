# Map Design

Map is designed using <a href="https://www.mapeditor.org/">Tiled Map Editor</a> and exported as a JSON file.

- Bottom-first layer, `Interactions` is used to define walkable areas and interactions with objects.
  - Rows 1 and 2 contain the walkable map elements
  - Rows 3 and 4 contain the obsacles
  - Rows 5 and 6 contain the collectible items
- All the layers start with `Player` contains various elements related to the character design.
  - Each row is used for separate level.
- All other layers are used for the background design.
