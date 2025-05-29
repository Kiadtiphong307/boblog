
export interface Article {
    id: number
    title: string
    slug: string
    content: string
    created_at: string
    createdAt: string
    author?: {
      username: string
    }
    category?: {
      id: number
      name: string
    }
    tags: {
      id: number
      name: string
    }[]
  }
  
  export interface Category {
    id: number
    name: string
  }
  
  export interface Comment {
    id: number
    content: string
    created_at: string
    updated_at: string
    user?: {
      username: string
    }
    
  }
  