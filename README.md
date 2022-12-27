# gOrder

A simple script written in Go to order the folder where it is. You have two options:

1. You can run it and read at runtime the json file: \
   `b, err := os.ReadFile("formats.json")`
2. You can embed the .json and just run the program without dependencies: \
  `//go:embed formats.json` \
  `var embedJson embed.FS` \
  `b, err := embedJson.ReadFile("formats.json")`
  