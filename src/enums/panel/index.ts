export enum PanelStateNetworkModeEnum {
  // 局域网
  'lan' = 0,
  // 互联网
  'wan' = 1,
}

export enum PanelPanelConfigStyleEnum {
  'icon' = 1, // 图标风格
  'info' = 0, // 详情风格
  'small' = 1, // 同icon
}

// 卡片类型
export enum PanelItemCardTypeEnum {
  'normal' = 1, // 普通图标
  'docker' = 2, // docker（未实现）
  'builtInApp' = 3, // 内置应用
}
