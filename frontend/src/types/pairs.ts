import type { ModDependency } from "./minecraft-mod";

export interface PairDepend {
	configDepend: ModDependency;
	fileDepend: string;
}
