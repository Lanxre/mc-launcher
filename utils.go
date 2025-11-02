package main

import (
	"slices"
	"cmp"
	"strings"
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


func (s *FuncService) SortByLoader(mods []MinecraftMod, loader string) []MinecraftMod {
    var filteredMods []MinecraftMod
    
    for _, mod := range mods {
        for _, detail := range mod.Details {
            if detail.Loader == loader {
                filteredMods = append(filteredMods, mod)
                break
            }
            
            if strings.Contains(detail.Loader, ",") {
                loaders := strings.SplitSeq(detail.Loader, ",")
                for l := range loaders {
                    if strings.TrimSpace(l) == loader {
                        filteredMods = append(filteredMods, mod)
                        break
                    }
                }
            }
        }
    }
    
    return filteredMods
}
