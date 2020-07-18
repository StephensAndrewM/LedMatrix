package main

import(
    log "github.com/sirupsen/logrus"
)

type Icon struct {
    Name string
    Width int
    Height int
    Layout [][]uint8
}

var iconSet map[string]Icon

func RegisterIcon(name string, layout [][]uint8) {
    icon := Icon{}
    icon.Name = name
    icon.Width = len(layout[0])
    icon.Height = len(layout)
    icon.Layout = layout
    iconSet[name] = icon
}

func InitIcons() {

    // Initialize the map
    iconSet = make(map[string]Icon)

    RegisterIcon("biohazard-16", [][]uint8{
        {0,1,1,0,0,0,0,0,0,0,0,0,1,1,0,0},
        {1,1,1,0,0,0,0,0,0,0,0,0,1,1,1,0},
        {1,1,1,1,0,0,0,0,0,0,0,1,1,1,1,0},
        {0,0,1,0,0,1,1,1,1,1,0,0,1,0,0,0},
        {0,0,0,0,1,1,1,1,1,1,1,0,0,0,0,0},
        {0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,0},
        {0,0,0,1,1,0,0,1,0,0,1,1,0,0,0,0},
        {0,0,0,1,0,0,0,1,0,0,0,1,0,0,0,0},
        {0,0,0,1,1,0,1,1,1,0,1,1,0,0,0,0},
        {0,0,0,0,1,1,1,0,1,1,1,0,0,0,0,0},
        {0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,0},
        {0,0,0,0,0,1,1,1,1,1,0,0,0,0,0,0},
        {0,0,1,0,0,1,0,1,0,1,0,0,1,0,0,0},
        {1,1,1,1,0,0,0,0,0,0,0,1,1,1,1,0},
        {1,1,1,0,0,0,0,0,0,0,0,0,1,1,1,0},
        {0,1,1,0,0,0,0,0,0,0,0,0,1,1,0,0},
    })

    RegisterIcon("house-16", [][]uint8{
        {0,0,0,0,0,0,0,1,0,0,0,1,1,0,0,0},
        {0,0,0,0,0,0,1,1,1,0,0,1,1,0,0,0},
        {0,0,0,0,0,1,1,0,1,1,0,1,1,0,0,0},
        {0,0,0,0,1,1,0,1,0,1,1,1,1,0,0,0},
        {0,0,0,1,1,0,1,1,1,0,1,1,1,0,0,0},
        {0,0,1,1,0,1,1,1,1,1,0,1,1,0,0,0},
        {0,1,1,0,1,1,1,1,1,1,1,0,1,1,0,0},
        {1,1,0,1,1,1,1,1,1,1,1,1,0,1,1,0},
        {1,0,1,1,1,1,1,1,1,1,1,1,1,0,1,0},
        {0,0,1,0,0,1,1,1,1,1,0,0,1,0,0,0},
        {0,0,1,0,0,1,0,0,0,1,0,0,1,0,0,0},
        {0,0,1,1,1,1,0,0,0,1,1,1,1,0,0,0},
        {0,0,1,1,1,1,0,0,0,1,1,1,1,0,0,0},
        {0,0,1,1,1,1,0,0,0,1,1,1,1,0,0,0},
        {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
    })

    RegisterIcon("missing", [][]uint8{
        {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
        {1,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1},
        {1,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1},
        {1,0,0,0,0,0,0,0,0,0,0,0,1,0,0,1},
        {1,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1},
        {1,0,0,0,0,0,0,0,0,0,1,0,0,0,0,1},
        {1,0,0,0,0,0,0,0,0,1,0,0,0,0,0,1},
        {1,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1},
        {1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,1},
        {1,0,0,0,0,0,1,0,0,0,0,0,0,0,0,1},
        {1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1},
        {1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,1},
        {1,0,0,1,0,0,0,0,0,0,0,0,0,0,0,1},
        {1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,1},
        {1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
        {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
    })

}

func GetIcon(name string) Icon {
    icon, ok := iconSet[name]
    if !ok {
        icon, ok = iconSet["missing"]
        if !ok {
            log.Error("Could not load fallback icon.")
        }
    }
    return icon
}