package main

import (
	"slices"
	"cmp"
)

type FuncService struct {}

func NewFuncService() *FuncService {
    return &FuncService{}
}

func (s *FuncService) SortByVersions(mods []MinecraftMod, version string) []MinecraftMod {
    var filteredMods []MinecraftMod
    
    for _, mod := range mods {
        for _, v := range mod.Versions {
            if v == version {
                slices.SortFunc(mod.Versions, func(a, b string) int {
                    return cmp.Compare(a, b)
                })
                filteredMods = append(filteredMods, mod)
                break
            }
        }
    }
    
    return filteredMods
}


func (s *FuncService) SortByLoader(mods []MinecraftMod, searchedLoader string) []MinecraftMod {
    var filteredMods []MinecraftMod
    
    for _, mod := range mods {
        if slices.Contains(mod.Loaders, searchedLoader) {
            filteredMods = append(filteredMods, mod)
        }
    }
    
    return filteredMods
}
