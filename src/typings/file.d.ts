declare namespace File {

	interface Info extends Common.InfoBase {
		src: string
		userId: number
		fileName: string
		method: number
		ext: string
		type?: number
		isPublicGallery?: boolean
		category?: string
		size?: number
	}

	// 图库分类
	type Category = 'icon' | 'wallpaper' | 'unknown'

}
