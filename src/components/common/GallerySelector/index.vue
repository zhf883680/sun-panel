<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { NButton, NEmpty, NGrid, NGridItem, NImage, NInput, NSelect, NSpin, NUpload, useMessage } from 'naive-ui'
import { getHdIconsImageUrl, getHdIconsList } from '@/api/system/gallery'
import type { HdIconItem, HdIconsListData } from '@/api/system/gallery'
import { getList as getFileList } from '@/api/system/file'
import { useAuthStore } from '@/store'
import { t } from '@/locales'

interface Props {
  show: boolean
  defaultType?: string
}

const props = withDefaults(defineProps<Props>(), {
  defaultType: 'icon',
})
const emit = defineEmits<{
  (e: 'update:show', show: boolean): void
  (e: 'select', src: string): void
}>()

const ms = useMessage()
const authStore = useAuthStore()

// 分页加载数量（本地与在线共用滚动加载）
const PAGE_SIZE = 60

const tab = ref<'local' | 'online'>('local')
const fileType = ref<string>(props.defaultType)
const keyword = ref('')
const localList = ref<File.Info[]>([])
const localLoading = ref(false)
const page = ref(1)

// 在线图库：数据来自后端本地缓存，服务端分页 + 搜索
const onlineList = ref<HdIconItem[]>([])
const onlineLoading = ref(false)
const onlineSyncing = ref(false)
const onlineTotal = ref(0)
const onlinePage = ref(0)
const onlineStyle = ref('all')

const typeOptions = computed(() => [
  { label: t('gallery.typeIcon'), value: 'icon' },
  { label: t('gallery.typeWallpaper'), value: 'wallpaper' },
])

// 在线图库风格筛选（服务端过滤）
const styleOptions = [
  { label: t('gallery.styleAll'), value: 'all' },
  { label: t('gallery.styleBorderRadius'), value: 'border-radius' },
  { label: t('gallery.styleCircle'), value: 'circle' },
  { label: t('gallery.styleSvg'), value: 'svg' },
]

const filterLocalList = computed(() => {
  const key = keyword.value.trim().toLowerCase()
  if (!key)
    return localList.value
  return localList.value.filter(item => item.fileName?.toLowerCase().includes(key))
})

// 本地分页（滚动加载更多）
const showLocalList = computed(() => filterLocalList.value.slice(0, page.value * PAGE_SIZE))
const localHasMore = computed(() => showLocalList.value.length < filterLocalList.value.length)

// 在线分页（服务端已分页，直接展示）
const onlineHasMore = computed(() => onlineList.value.length < onlineTotal.value)

watch(() => props.show, (show) => {
  if (show) {
    fileType.value = props.defaultType || 'icon'
    keyword.value = ''
    page.value = 1
    if (tab.value === 'local')
      getLocalList()
    else
      loadOnlinePage(true)
  }
})

watch(tab, (value) => {
  if (value === 'local')
    getLocalList()
  else
    loadOnlinePage(true)
})

watch(fileType, () => {
  if (tab.value === 'local')
    getLocalList()
})

watch(onlineStyle, () => {
  if (tab.value === 'online')
    loadOnlinePage(true)
})

watch(keyword, () => {
  page.value = 1
  if (tab.value === 'online')
    scheduleOnlineSearch()
})

let searchTimer: ReturnType<typeof setTimeout> | undefined
let retryTimer: ReturnType<typeof setTimeout> | undefined

function scheduleOnlineSearch() {
  if (searchTimer)
    clearTimeout(searchTimer)
  searchTimer = setTimeout(() => loadOnlinePage(true), 300)
}

// 首次同步尚未完成时轮询重试，同步完成后自动展示
function scheduleOnlineRetry() {
  if (retryTimer)
    clearTimeout(retryTimer)
  retryTimer = setTimeout(() => {
    if (tab.value === 'online')
      loadOnlinePage(true)
  }, 3000)
}

async function getLocalList() {
  localLoading.value = true
  try {
    const { data } = await getFileList<Common.ListResponse<File.Info[]>>(fileType.value)
    localList.value = data?.list || []
  }
  catch {
    localList.value = []
  }
  localLoading.value = false
}

async function loadOnlinePage(reset = false) {
  if (onlineLoading.value)
    return
  if (reset) {
    onlinePage.value = 0
    onlineTotal.value = 0
    onlineList.value = []
  }
  else if (onlineTotal.value > 0 && onlineList.value.length >= onlineTotal.value) {
    return
  }

  onlineLoading.value = true
  try {
    const nextPage = onlinePage.value + 1
    const { code, data } = await getHdIconsList<HdIconsListData>({
      page: nextPage,
      pageSize: PAGE_SIZE,
      keyword: keyword.value.trim(),
      style: onlineStyle.value,
    })
    if (code === 0 && data) {
      onlineSyncing.value = !!data.syncing
      onlineTotal.value = data.total || 0
      onlinePage.value = nextPage
      const items = data.list || []
      onlineList.value = reset ? items : onlineList.value.concat(items)
      // 同步中且暂无数据：稍后自动重试
      if (onlineSyncing.value && onlineTotal.value === 0)
        scheduleOnlineRetry()
    }
    else {
      ms.error(t('gallery.onlineLoadFailed'))
    }
  }
  catch {
    ms.error(t('gallery.onlineLoadFailed'))
  }
  onlineLoading.value = false
}

// 滚动到底部加载更多
function handleScroll(e: Event) {
  const target = e.target as HTMLElement
  if (target.scrollTop + target.clientHeight < target.scrollHeight - 40)
    return
  if (tab.value === 'local') {
    if (localHasMore.value)
      page.value++
  }
  else if (onlineHasMore.value) {
    loadOnlinePage(false)
  }
}

// 选择本地图库图片
function handleSelectLocal(item: File.Info) {
  emit('select', item.src)
  emit('update:show', false)
}

// 选择在线图库图标
function handleSelectOnline(item: HdIconItem) {
  emit('select', getHdIconsImageUrl(item.style, item.file))
  emit('update:show', false)
}

function handleUploadFinish() {
  getLocalList()
}
</script>

<template>
  <div v-if="show" class="gallery-overlay" @click.self="emit('update:show', false)">
    <div class="gallery-panel">
      <div class="flex items-center justify-between mb-[10px]">
        <div class="text-base font-bold">
          {{ $t('gallery.title') }}
        </div>
        <NButton size="tiny" quaternary @click="emit('update:show', false)">
          ✕
        </NButton>
      </div>

      <!-- Tab -->
      <div class="gallery-tabs">
        <div class="gallery-tab" :class="{ active: tab === 'local' }" @click="tab = 'local'">
          {{ $t('gallery.localGallery') }}
        </div>
        <div class="gallery-tab" :class="{ active: tab === 'online' }" @click="tab = 'online'">
          {{ $t('gallery.onlineGallery') }}
        </div>
      </div>

      <!-- 工具栏 -->
      <div class="flex items-center gap-[8px] my-[10px] flex-wrap">
        <NInput
          v-model:value="keyword" size="small" clearable
          :placeholder="$t('gallery.searchPlaceholder')" style="flex:1;min-width:160px;"
        />
        <NSelect v-if="tab === 'local'" v-model:value="fileType" size="small" :options="typeOptions" style="width:130px;" />
        <NSelect v-if="tab === 'online'" v-model:value="onlineStyle" size="small" :options="styleOptions" style="width:150px;" />
        <div v-if="tab === 'local'" class="flex items-center">
          <NUpload
            action="/api/file/uploadImg"
            name="imgfile"
            :show-file-list="false"
            :data="{ category: fileType }"
            :headers="{ token: authStore.token as string }"
            accept="image/*"
            @finish="handleUploadFinish"
          >
            <NButton size="small" type="primary">
              {{ $t('gallery.upload') }}
            </NButton>
          </NUpload>
        </div>
      </div>

      <!-- 内容 -->
      <div class="gallery-body" @scroll="handleScroll">
        <NSpin :show="tab === 'local' ? localLoading : onlineLoading">
          <!-- 本地图库 -->
          <template v-if="tab === 'local'">
            <NEmpty v-if="filterLocalList.length === 0 && !localLoading" :description="$t('gallery.empty')" class="mt-[40px]" />
            <NGrid v-else cols="3 400:4 600:5 800:6" :x-gap="8" :y-gap="8">
              <NGridItem v-for="item in showLocalList" :key="item.id">
                <div class="gallery-item" @click="handleSelectLocal(item)">
                  <NImage :src="item.src" lazy object-fit="contain" class="gallery-img" />
                  <div class="gallery-name">
                    {{ item.fileName }}
                  </div>
                </div>
              </NGridItem>
            </NGrid>
            <div v-if="localHasMore" class="text-center text-xs opacity-60 py-[10px]">
              {{ $t('gallery.scrollForMore') }}
            </div>
          </template>

          <!-- 在线图库 -->
          <template v-else>
            <div v-if="onlineSyncing && onlineList.length === 0" class="text-center text-xs opacity-70 py-[40px]">
              {{ $t('gallery.onlineSyncing') }}
            </div>
            <NEmpty v-else-if="onlineList.length === 0 && !onlineLoading" :description="$t('gallery.empty')" class="mt-[40px]" />
            <NGrid v-else cols="4 400:5 600:6 800:8" :x-gap="8" :y-gap="8">
              <NGridItem v-for="item in onlineList" :key="`${item.style}/${item.file}`">
                <div class="gallery-item" :title="item.name" @click="handleSelectOnline(item)">
                  <img :src="getHdIconsImageUrl(item.style, item.file)" loading="lazy" class="gallery-img" alt="">
                  <div class="gallery-name">
                    {{ item.name }}
                  </div>
                </div>
              </NGridItem>
            </NGrid>
            <div v-if="onlineHasMore" class="text-center text-xs opacity-60 py-[10px]">
              {{ $t('gallery.scrollForMore') }}
            </div>
          </template>
        </NSpin>
      </div>

      <div class="text-xs opacity-50 mt-[8px] text-center">
        {{ tab === 'online' ? $t('gallery.onlineTip') : '' }}
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 与 v1.8.1 保持一致：图库网格不铺透明棋盘格（棋盘格只用于图标预览与上传管理） */
.gallery-overlay {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 3000;
  background-color: #00000073;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gallery-panel {
  width: 860px;
  max-width: 94vw;
  height: 640px;
  max-height: 88vh;
  background-color: var(--n-color, #fff);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 40px #0000004d;
}

.gallery-body {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;
}

.gallery-tabs {
  display: flex;
  border-bottom: 1px solid rgba(128, 128, 128, 0.2);
  margin-bottom: 8px;
}

.gallery-tab {
  padding: 6px 14px;
  cursor: pointer;
  font-size: 13px;
  border-bottom: 2px solid transparent;
  transition: color 0.2s, border-color 0.2s;
  opacity: 0.7;
}

.gallery-tab:hover {
  opacity: 1;
}

.gallery-tab.active {
  border-bottom-color: #18a058;
  color: #18a058;
  opacity: 1;
  font-weight: 600;
}

.gallery-item {
  cursor: pointer;
  border-radius: 8px;
  padding: 6px;
  transition: background-color 0.2s, transform 0.2s;
  border: 1px solid transparent;
}

.gallery-item:hover {
  background-color: #80808026;
  transform: translateY(-2px);
}

.gallery-img {
  width: 100%;
  height: 64px;
  object-fit: contain;
}

.gallery-name {
  font-size: 11px;
  margin-top: 4px;
  text-align: center;
  opacity: 0.75;
  white-space: nowrap;
  overflow: hidden;
}
</style>
