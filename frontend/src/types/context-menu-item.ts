export interface ContextMenuItem {
  id: string
  label: string
  icon?: string
  shortcut?: string
  disabled?: boolean
  danger?: boolean
  children?: ContextMenuItem[]
  handler?: () => void
}