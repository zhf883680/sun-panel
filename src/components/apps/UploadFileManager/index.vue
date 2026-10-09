<script setup lang="ts">
import { NAlert, NButton, NButtonGroup, NCard, NEllipsis, NEmpty, NGrid, NGridItem, NImage, NImageGroup, NInput, NSelect, NSpin, NTabPane, NTabs, NUpload, useDialog, useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { computed, onMounted, ref, watch } from 'vue'
import { deletes, getList, updateCategory } from '@/api/system/file'
import { getHdIconsIndex, getHdIconsImageUrl } from '@/api/system/gallery'
import type { HdIconItem, HdIconsIndex } from '@/api/system/gallery'
import { set as savePanelConfig } from '@/api/panel/userConfig'
import { RoundCardModal, SvgIcon } from '@/components/common'
import { copyToClipboard, timeFormat } from '@/utils/cmn'
import { t } from '@/locales'
import { useAuthStore, usePanelState } from '@/store'

interface InfoModalState {
  title: string
  show: boolean
  fileInfo: File.Info | null
}

// 在线图库分页加载数量
const PAGE_SIZE = 60

const imageList = ref<File.Info[]>([])
const ms = useMessage()
const dialog = useDialog()
const panelStore = usePanelState()
const authStore = useAuthStore()
const loading = ref(false)
const currentTab = ref<'local' | 'online'>('local')
const category = ref<'all' | 'icon' | 'wallpaper'>('all')
const keyword = ref('')
const onlineList = ref<HdIconItem[]>([])
const onlineLoading = ref(false)
const onlineLoaded = ref(false)
const onlineStyle = ref('border-radius')
const page = ref(1)
const infoModalState = ref<InfoModalState>({
  show: false,
  title: '',
  fileInfo: null,
})

const categoryOptions = computed(() => [
  { label: t('apps.uploadsFileManager.categoryAll'), value: 'all' },
  { label: t('apps.uploadsFileManager.categoryIcon'), value: 'icon' },
  { label: t('apps.uploadsFileManager.categoryWallpaper'), value: 'wallpaper' },
])

// 与 v1.8.1 保持一致：该下拉仅作展示，图标使用自身风格地址
const styleOptions = [
  { label: t('gallery.styleBorderRadius'), value: 'border-radius' },
  { label: t('gallery.styleCircle'), value: 'circle' },
  { label: t('gallery.styleSvg'), value: 'svg' },
]

const filterImageList = computed(() => {
  const key = keyword.value.trim().toLowerCase()
  let list = imageList.value
  if (category.value !== 'all') {
    // 与后端一致：历史上传的图片没有分类，默认按图标处理
    list = list.filter(item => (item.category || 'icon') === category.value)
  }
  if (key)
    list = list.filter(item => item.fileName?.toLowerCase().includes(key))
  return list
})

// 本地图库分页（滚动加载更多）
const showLocalList = computed(() => filterImageList.value.slice(0, page.value * PAGE_SIZE))
const localHasMore = computed(() => showLocalList.value.length < filterImageList.value.length)

const filterOnlineList = computed(() => {
  const key = keyword.value.trim().toLowerCase()
  if (!key)
    return onlineList.value
  return onlineList.value.filter(item => item.name.toLowerCase().includes(key))
})

const showOnlineList = computed(() => filterOnlineList.value.slice(0, page.value * PAGE_SIZE))
const hasMore = computed(() => showOnlineList.value.length < filterOnlineList.value.length)

async function getFileList() {
  loading.value = true
  try {
    const { data } = await getList<Common.ListResponse<File.Info[]>>()
    imageList.value = data?.list || []
  }
  catch {
    imageList.value = []
  }
  loading.value = false
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

watch(currentTab, (value) => {
  if (value === 'online' && !onlineLoaded.value)
    getOnlineList()
})

watch([keyword, category, currentTab], () => {
  page.value = 1
})

function handleScroll(e: Event) {
  const target = e.target as HTMLElement
  if (target.scrollTop + target.clientHeight < target.scrollHeight - 40)
    return
  if (currentTab.value === 'local' ? localHasMore.value : hasMore.value)
    page.value++
}

async function copyImageUrl(text: string) {
  const res = await copyToClipboard(text)
  if (res)
    ms.success(t('apps.uploadsFileManager.copySuccess'))

  else
    ms.error(t('apps.uploadsFileManager.copyFailed'))
}

function handleDelete(id: number) {
  dialog.warning({
    title: t('common.warning'),
    content: t('apps.uploadsFileManager.deleteWarningText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      deletesImges(id)
    },
  })
}

async function deletesImges(id: number) {
  try {
    const { code, msg } = await deletes([id])
    if (code === 0) {
      getFileList()
      ms.success(t('common.success'))
    }
    else {
      ms.error(`${t('common.failed')}:${msg}`)
    }
  }
  catch (error) {
    ms.error(t('common.failed'))
  }
}

function handleInfoClick(fileInfo: File.Info) {
  infoModalState.value.fileInfo = fileInfo
  infoModalState.value.show = true
}

function handleSetWallpaper(imgSrc: string) {
  panelStore.panelConfig.backgroundImageSrc = imgSrc
  savePanelConfig({ panel: panelStore.panelConfig })
}

// 图标 / 壁纸 分类切换
async function handleToggleCategory(fileInfo: File.Info) {
  const nextCategory: File.Category = fileInfo.category === 'wallpaper' ? 'icon' : 'wallpaper'
  try {
    const { code, msg } = await updateCategory(fileInfo.id as number, nextCategory)
    if (code === 0) {
      fileInfo.category = nextCategory
      ms.success(t('common.success'))
    }
    else {
      ms.error(`${t('common.failed')}:${msg}`)
    }
  }
  catch (error) {
    ms.error(t('common.failed'))
  }
}

function handleUploadFinish({ event }: { file: UploadFileInfo, event?: ProgressEvent }) {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  if (res.code === 0)
    getFileList()
  else
    ms.error(`${t('common.failed')}:${res.msg}`)
}

onMounted(() => {
  getFileList()
})
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 p-2 h-full flex flex-col">
    <NAlert type="info" :bordered="false">
      {{ $t('apps.uploadsFileManager.alertText') }}
    </NAlert>

    <NTabs v-model:value="currentTab" type="line" size="small" class="mt-1">
      <!-- naive-ui 2.45 起，NTabPane 无子节点时标签会被空插槽覆盖，这里补上 #tab 插槽 -->
      <NTabPane name="local" :tab="$t('gallery.localGallery')">
        <template #tab>{{ $t('gallery.localGallery') }}</template>
      </NTabPane>
      <NTabPane name="online" :tab="$t('gallery.onlineGallery')">
        <template #tab>{{ $t('gallery.onlineGallery') }}</template>
      </NTabPane>
    </NTabs>

    <div class="flex items-center gap-[8px] my-[10px] flex-wrap">
      <NInput
        v-model:value="keyword" size="small" clearable
        :placeholder="$t('gallery.searchPlaceholder')" style="flex:1;min-width:160px;"
      />
      <NSelect v-if="currentTab === 'local'" v-model:value="category" size="small" :options="categoryOptions" style="width:130px;" />
      <NSelect v-if="currentTab === 'online'" v-model:value="onlineStyle" size="small" :options="styleOptions" style="width:150px;" />
      <div v-if="currentTab === 'local'" class="flex items-center">
        <NUpload
          action="/api/file/uploadImg"
          name="imgfile"
          :show-file-list="false"
          :data="{ category: category === 'all' ? 'icon' : category }"
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

    <div class="flex-1 overflow-y-auto" @scroll="handleScroll">
      <NSpin :show="currentTab === 'local' ? loading : onlineLoading">
        <!-- 本地图库 -->
        <template v-if="currentTab === 'local'">
          <div v-if="filterImageList.length === 0 && !loading" class="flex justify-center">
            {{ $t('apps.uploadsFileManager.nothingText') }}
          </div>
          <NImageGroup v-else>
            <NGrid cols="2 300:2 600:4 900:6 1100:9" :x-gap="5" :y-gap="5">
              <NGridItem v-for=" item, index in showLocalList" :key="index">
                <NCard size="small" style="border-radius: 5px;" :bordered="true">
                  <template #cover>
                    <div class="card transparent-grid">
                      <NImage :lazy="true" style="object-fit: contain;height: 100%;" :src="item.src" />
                    </div>
                  </template>
                  <template #footer>
                    <span class="text-xs">
                      <NEllipsis>
                        {{ item.fileName }}
                      </NEllipsis>
                    </span>
                    <div class="flex justify-center mt-[10px]">
                      <NButtonGroup>
                        <NButton size="tiny" tertiary style="cursor: pointer;" :title="$t('apps.uploadsFileManager.copyLink')" @click="copyImageUrl(item.src)">
                          <template #icon>
                            <SvgIcon icon="ion-copy" />
                          </template>
                        </NButton>
                        <NButton size="tiny" tertiary style="cursor: pointer;" :title="timeFormat(item.createTime)" @click="handleInfoClick(item)">
                          <template #icon>
                            <SvgIcon icon="mdi-information-box-outline" />
                          </template>
                        </NButton>
                        <NButton size="tiny" tertiary style="cursor: pointer;" :title="$t('apps.uploadsFileManager.setWallpaper')" @click="handleSetWallpaper(item.src)">
                          <template #icon>
                            <SvgIcon icon="lucide:wallpaper" />
                          </template>
                        </NButton>
                        <NButton
                          size="tiny" tertiary style="cursor: pointer;"
                          :title="item.category === 'wallpaper' ? $t('apps.uploadsFileManager.markAsIcon') : $t('apps.uploadsFileManager.markAsWallpaper')"
                          @click="handleToggleCategory(item)"
                        >
                          <template #icon>
                            <SvgIcon icon="lucide:tag" />
                          </template>
                        </NButton>
                        <NButton size="tiny" tertiary type="error" style="cursor: pointer;" :title="$t('common.delete')" @click="handleDelete(item.id as number)">
                          <template #icon>
                            <SvgIcon icon="material-symbols-delete" />
                          </template>
                        </NButton>
                      </NButtonGroup>
                    </div>
                  </template>
                </NCard>
              </NGridItem>
            </NGrid>
          </NImageGroup>
          <div v-if="localHasMore" class="text-center text-xs opacity-60 py-[10px]">
            {{ $t('gallery.scrollForMore') }}
          </div>
        </template>

        <!-- 在线图库 -->
        <template v-else>
          <NEmpty v-if="showOnlineList.length === 0 && !onlineLoading" :description="$t('gallery.empty')" class="mt-[40px]" />
          <NGrid v-else cols="4 400:5 600:6 800:8 1100:10" :x-gap="5" :y-gap="5">
            <NGridItem v-for="item in showOnlineList" :key="item.name">
              <NCard size="small" style="border-radius: 5px;" :bordered="true">
                <template #cover>
                  <div class="card transparent-grid">
                    <img
                      :src="getHdIconsImageUrl(item.style, item.file)"
                      loading="lazy" style="object-fit: contain;height: 100%;width: 100%;" alt=""
                    >
                  </div>
                </template>
                <template #footer>
                  <span class="text-xs">
                    <NEllipsis>
                      {{ item.name }}
                    </NEllipsis>
                  </span>
                  <div class="flex justify-center mt-[10px]">
                    <NButtonGroup>
                      <NButton size="tiny" tertiary style="cursor: pointer;" :title="$t('apps.uploadsFileManager.copyLink')" @click="copyImageUrl(getHdIconsImageUrl(item.style, item.file))">
                        <template #icon>
                          <SvgIcon icon="ion-copy" />
                        </template>
                      </NButton>
                      <NButton size="tiny" tertiary style="cursor: pointer;" :title="$t('apps.uploadsFileManager.setWallpaper')" @click="handleSetWallpaper(getHdIconsImageUrl(item.style, item.file))">
                        <template #icon>
                          <SvgIcon icon="lucide:wallpaper" />
                        </template>
                      </NButton>
                    </NButtonGroup>
                  </div>
                </template>
              </NCard>
            </NGridItem>
          </NGrid>
          <div v-if="hasMore" class="text-center text-xs opacity-60 py-[10px]">
            {{ $t('gallery.scrollForMore') }}
          </div>
        </template>
      </NSpin>
    </div>

    <RoundCardModal v-model:show="infoModalState.show" style="max-width: 300px;" size="small" :title="$t('apps.uploadsFileManager.infoTitle')">
      <div>
        <div>
          <div class="mb-2">
            <span class="text-slate-500">
              {{ $t('apps.uploadsFileManager.fileName') }}
            </span>
            <div class="text-xs">
              {{ infoModalState.fileInfo?.fileName }}
            </div>
          </div>
          <div class="mb-2">
            <span class="text-slate-500">
              {{ $t('apps.uploadsFileManager.path') }}
            </span>
            <div class="text-xs">
              {{ infoModalState.fileInfo?.src }}
            </div>
          </div>
          <div class="mb-2">
            <span class="text-slate-500">
              {{ $t('apps.uploadsFileManager.uploadTime') }}
            </span>
            <div class="text-xs">
              {{ timeFormat(infoModalState.fileInfo?.createTime) }}
            </div>
          </div>
        </div>
      </div>
    </RoundCardModal>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80px;
  overflow: hidden;
}

.transparent-grid {
  background-image: linear-gradient(45deg, #f0f0f0 25%, transparent 25%, transparent 75%, #f0f0f0 75%),
    linear-gradient(45deg, #f0f0f0 25%, transparent 25%, transparent 75%, #f0f0f0 75%);
  background-size: 16px 16px;
  background-position: 0 0, 8px 8px;
}
</style>
