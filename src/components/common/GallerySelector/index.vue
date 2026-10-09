<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { NButton, NEmpty, NGrid, NGridItem, NImage, NInput, NSelect, NSpin, NUpload, useMessage } from 'naive-ui'
import { getHdIconsIndex, getHdIconsImageUrl } from '@/api/system/gallery'
import type { HdIconItem, HdIconsIndex } from '@/api/system/gallery'
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

// 分页加载数量
const PAGE_SIZE = 60

const tab = ref<'local' | 'online'>('local')
const fileType = ref<string>(props.defaultType)
const keyword = ref('')
const localList = ref<File.Info[]>([])
const localLoading = ref(false)
const onlineList = ref<HdIconItem[]>([])
const onlineLoading = ref(false)
const onlineLoaded = ref(false)
const onlineStyle = ref('border-radius')
const page = ref(1)

const typeOptions = computed(() => [
  { label: t('gallery.typeIcon'), value: 'icon' },
  { label: t('gallery.typeWallpaper'), value: 'wallpaper' },
])

// 与 v1.8.1 保持一致：该下拉仅作展示，图标使用自身风格地址
const styleOptions = [
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

const filterOnlineList = computed(() => {
  const key = keyword.value.trim().toLowerCase()
  if (!key)
    return onlineList.value
  return onlineList.value.filter(item => item.name.toLowerCase().includes(key))
})

// 分页（滚动加载更多）
const showLocalList = computed(() => filterLocalList.value.slice(0, page.value * PAGE_SIZE))
const localHasMore = computed(() => showLocalList.value.length < filterLocalList.value.length)
const showOnlineList = computed(() => filterOnlineList.value.slice(0, page.value * PAGE_SIZE))
const hasMore = computed(() => showOnlineList.value.length < filterOnlineList.value.length)

watch(() => props.show, (show) => {
  if (show) {
    fileType.value = props.defaultType || 'icon'
    keyword.value = ''
    page.value = 1
    if (tab.value === 'local')
      getLocalList()
  }
})

watch(tab, (value) => {
  if (value === 'local')
    getLocalList()
  else if (!onlineLoaded.value)
    getOnlineList()
})

watch(fileType, () => {
  if (tab.value === 'local')
    getLocalList()
})

watch(keyword, () => {
  page.value = 1
})

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

async function getOnlineList() {
  onlineLoading.value = true
  try {
    const { code, data } = await getHdIconsIndex<HdIconsIndex>()
    if (code === 0 && data?.list) {
      onlineList.value = data.list
      onlineLoaded.value = true
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
  if (tab.value === 'local' ? localHasMore.value : hasMore.value)
    page.value++
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
            <NEmpty v-if="showOnlineList.length === 0 && !onlineLoading" :description="$t('gallery.empty')" class="mt-[40px]" />
            <NGrid v-else cols="4 400:5 600:6 800:8" :x-gap="8" :y-gap="8">
              <NGridItem v-for="item in showOnlineList" :key="item.name">
                <div class="gallery-item" :title="item.name" @click="handleSelectOnline(item)">
                  <img :src="getHdIconsImageUrl(item.style, item.file)" loading="lazy" class="gallery-img" alt="">
                  <div class="gallery-name">
                    {{ item.name }}
                  </div>
                </div>
              </NGridItem>
            </NGrid>
            <div v-if="hasMore" class="text-center text-xs opacity-60 py-[10px]">
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
