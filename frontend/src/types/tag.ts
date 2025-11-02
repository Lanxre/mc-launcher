export interface Tag {
  id: string | number;
  name: string;
  imageUrl?: string;
  altText?: string;
  clickable?: boolean;
  closable?: boolean;
  size?: 'small' | 'medium' | 'large';
  variant?: 'default' | 'outlined' | 'ghost';
  rounded?: boolean;
}
