import { defineStore } from 'pinia'
import { MinecraftMod } from '../types'

export const useModStore = defineStore('mod', {
  state: () => ({
    currentMod: null as MinecraftMod | null,
    currentParsePage: null as  number | null,
    allMods: [] as MinecraftMod[],
    isLoading: false,
    error: null as string | null
  }),

  getters: {
    getCurrentMod: (state) => state.currentMod,
    getAllMods: (state) => state.allMods,
    getModsCount: (state) => state.allMods.length,
    getParsePage: (state) => state.currentParsePage,
    
    searchMods: (state) => (query: string) => {
      return state.allMods.filter(mod => 
        mod.Name.toLowerCase().includes(query.toLowerCase())
      )
    },
    
    getModsByLoader: (state) => (loader: string) => {
      return state.allMods.filter(mod => 
        mod.Details?.some(detail => 
          detail.Loader.toLowerCase().includes(loader.toLowerCase())
        )
      )
    }
  },

  actions: {
    setAllMods(mods: MinecraftMod[]) {
      this.allMods = mods
    },

    addMods(mods: MinecraftMod[]) {
      const existingLinks = new Set(this.allMods.map(mod => mod.ModPageLink))
      const uniqueNewMods = mods.filter(mod => !existingLinks.has(mod.ModPageLink))
      this.allMods.push(...uniqueNewMods)
    },

    setCurrentMod(mod: MinecraftMod) {
      this.currentMod = mod
    },

    setCurrentParsePage(page: number) {
      this.currentParsePage = page
    },

    clearCurrentMod() {
      this.currentMod = null
    },

    clearAllMods() {
      this.allMods = []
      this.currentMod = null
    },

    setLoading(loading: boolean) {
      this.isLoading = loading
    },

    setError(message: string | null) {
      this.error = message
    }
  }
})