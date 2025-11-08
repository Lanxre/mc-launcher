import type { MinecraftMod, ModDependency } from '@/types'
import { OpenExternalLink } from '../../../wailsjs/go/main/App'
import { SaveYamlModConfig } from '../../../wailsjs/go/main/FuncService'
import { GetModDetails } from '../../../wailsjs/go/parser/ScraperService'

export const uniqueBy = <T>(array: T[], keyFn: (item: T) => string): T[] => {
  const seen = new Set<string>()
  const result: T[] = []

  for (const item of array) {
    const key = keyFn(item)
    if (!seen.has(key)) {
      seen.add(key)
      result.push(item)
    }
  }

  return result
}


export const openLink = async (url: string) => {
    await OpenExternalLink(url)
}

export const saveModToYaml = async (mod: MinecraftMod, filename: string) => {
  try {
    if (mod !== null && mod !== undefined) { 
      await SaveYamlModConfig(mod, filename) 
    }
  } catch (err) {
    console.error('Ошибка сохранения:', err)
  }
}

export const getMinecraftDownloadFileName = (modName: string, versions: string[]): string => {
  const normalizedModName = modName.toLowerCase().replaceAll(' ', '_');
  const normalizedVersions = versions.join('_').toLowerCase();
  
  return `${normalizedModName}_${normalizedVersions}.jar`;
};

export async function enrichDependencies(deps: ModDependency[]) {
  for (const dep of deps) {
    if (dep.Details == null || dep.Details.length === 0) {
      const data = await GetModDetails(dep.ModPageLink)
      dep.Details = data.Details
    }
  }
}