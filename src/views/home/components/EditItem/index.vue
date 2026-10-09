<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCollapse, NCollapseItem, NColorPicker, NForm, NFormItem, NGrid, NGridItem, NInput, NInputGroup, NModal, NSelect, useMessage } from 'naive-ui'
import IconEditor from './IconEditor.vue'
import AppIcon from '../AppIcon/index.vue'
import { edit, getSiteFavicon } from '@/api/panel/itemIcon'
import { getList as getGroupList } from '@/api/panel/itemIconGroup'
import { defautSwatchesBackground } from '@/utils/defaultData'
import { usePanelState } from '@/store'
import { PanelItemCardTypeEnum } from '@/enums'
import { t } from '@/locales'

interface Props {
  visible: boolean
  itemInfo: Panel.Info | null
  itemGroupId?: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const ms = useMessage()
const panelState = usePanelState()
const submitLoading = ref(false)
const getIconLoading = ref([false, false])
const showMoreOptions = ref<string[]>([])
const itemIconGroupOptions = ref<{
  label: string
  value: number
}[]>([])

const defaultCardBackground = '#2a2a2a6b'

const restoreDefault: Panel.Info = {
  icon: null,
  title: '',
  url: '',
  lanUrl: '',
  description: '',
  openMethod: 2,
  backgroundColor: defaultCardBackground,
  cardType: PanelItemCardTypeEnum.normal,
  onlyName: '',
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', item: Panel.Info): void// 创建完成
}

const model = ref<Panel.Info>(props.itemInfo ? { ...props.itemInfo } : { ...restoreDefault })
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  title: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },
  url: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: t('form.required'),
  },
}

const options = [
  {
    default: true,
    label: t('iconItem.currentPageOpen'),
    value: 1,
  },
  {
    label: t('iconItem.newWindowOpen'),
    value: 2,
  },
  {
    label: t('iconItem.currentPageLayerOpen'),
    value: 3,
  },
]

// 卡片类型（v1.7+）
const cardTypeOptions = [
  {
    label: t('iconItem.ordinaryIcon'),
    value: PanelItemCardTypeEnum.normal,
  },
  {
    label: t('iconItem.builtInApp'),
    value: PanelItemCardTypeEnum.builtInApp,
  },
]

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

async function editApi() {
  submitLoading.value = true
  try {
    const { code, data, msg } = await edit<Panel.ItemInfo>(model.value)
    if (code === 0) {
      show.value = false
      model.value = { ...restoreDefault }

      emit('done', data)
    }
    else if (code === -2) {
      // 唯一标识已被占用
      ms.error(`${t('common.saveFail')}: ${t('iconItem.onlyNameExisted')}`)
    }
    else {
      ms.error(`${t('common.saveFail')}:${msg}`)
    }
  }
  catch (error) {
    ms.error(t('common.saveFail'))
  }
  submitLoading.value = false
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      editApi()
  })
}

function handleResetBackgroundColor() {
  model.value.backgroundColor = defaultCardBackground
}

async function getIconByUrl(url: string, loadingIndex: number) {
  getIconLoading.value[loadingIndex] = true
  try {
    const { code, data } = await getSiteFavicon<{ iconUrl: string, title?: string }>(url)
    if (code === 0) {
      model.value.icon = {
        itemType: 2,
        src: data.iconUrl,
      }
      // 同步填充网站标题
      if (data.title)
        model.value.title = data.title.slice(0, 20)
    }
    else if (data && data.title) {
      // 图标获取失败，但拿到了标题
      model.value.title = data.title.slice(0, 20)
      ms.warning(t('iconItem.getIconTitleOnly'))
    }
    else {
      ms.error(t('iconItem.geticonFail'))
    }
  }
  catch (error) {
    ms.error(t('iconItem.geticonFail'))
  }
  getIconLoading.value[loadingIndex] = false
}

watch(() => props.visible, (newValue) => {
  if (newValue === true) {
    model.value = props.itemInfo ? { ...props.itemInfo } : { ...restoreDefault }
    if (!model.value.backgroundColor)
      model.value.backgroundColor = defaultCardBackground
    if (!model.value.cardType)
      model.value.cardType = PanelItemCardTypeEnum.normal
    if (props.itemGroupId)
      model.value.itemIconGroupId = props.itemGroupId
  }

  getGroupListOptions()
})

function getGroupListOptions() {
  getGroupList<Common.ListResponse<Panel.ItemIconGroup[]>>().then(({ data, code, msg }) => {
    if (code === 0) {
      itemIconGroupOptions.value = []

      for (let i = 0; i < data.list.length; i++) {
        const element = data.list[i]
        if (i === 0 && !model.value.itemIconGroupId) {
          model.value.itemIconGroupId = element.id
          restoreDefault.itemIconGroupId = element.id
        }

        itemIconGroupOptions.value.push({
          value: element.id as number,
          label: element.title as string,
        })
      }
    }
    else {
      ms.error(`${t('iconItem.getGroupFail')}:${msg}`)
    }
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 620px;border-radius: 1rem;" :title="itemInfo ? t('iconItem.edit') : t('iconItem.add')">
    <div class="max-h-[600px] overflow-auto p-[5px]">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NGrid cols="2" :x-gap="10" item-responsive>
          <NGridItem span="2 500:1">
            <NFormItem path="title" :label="$t('common.title')">
              <NInput v-model:value="model.title" type="text" show-count :maxlength="20" />
            </NFormItem>
          </NGridItem>
          <NGridItem span="2 500:1">
            <NFormItem path="openMethod" :label="t('iconItem.openMethod')">
              <NSelect v-model:value="model.openMethod" :options="options" />
            </NFormItem>
          </NGridItem>
        </NGrid>

        <NFormItem path="icon" :label="$t('common.icon')">
          <IconEditor v-model:item-icon="model.icon" :card-background="model.backgroundColor || defaultCardBackground" />
        </NFormItem>
        <NFormItem path="url" :label="t('iconItem.url')">
          <NInputGroup>
            <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="http(s)://" />
            <NButton :disabled="!model.url" :loading="getIconLoading[0]" @click="getIconByUrl(model.url, 0)">
              {{ $t('iconItem.getIcon') }}
            </NButton>
          </NInputGroup>
        </NFormItem>
        <NFormItem path="lanUrl" :label="t('iconItem.lanUrl')">
          <NInputGroup>
            <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
            <NButton :disabled="!model.lanUrl" :loading="getIconLoading[1]" @click="getIconByUrl(model.lanUrl || '', 1)">
              {{ $t('iconItem.getIcon') }}
            </NButton>
          </NInputGroup>
        </NFormItem>
        <NFormItem path="description" :label="$t('common.description')">
          <NInput v-model:value="model.description" type="text" show-count :maxlength="100" />
        </NFormItem>

        <!-- 更多选项 -->
        <div class="mt-[10px]">
          <NCollapse v-model:expanded-names="showMoreOptions">
            <NCollapseItem name="more">
              <template #header>
                <span class="font-bold">{{ t('iconItem.moreOptions') }}</span>
                <span class="ml-1 text-xs opacity-60">
                  ({{ t('iconItem.cardBackground') }} , {{ t('iconItem.iconGroup') }} , {{ t('iconItem.cardType') }} , {{ t('common.onlyName') }})
                </span>
              </template>

              <NFormItem :label="t('iconItem.cardBackground')">
                <div class="w-full">
                  <div class="flex items-center gap-[10px] flex-wrap">
                    <div class="card-bg-color-picker">
                      <NColorPicker
                        v-model:value="model.backgroundColor"
                        size="small"
                        style="width: 100%;"
                        :modes="['hex']"
                        :show-alpha="true"
                        :swatches="defautSwatchesBackground"
                      />
                    </div>
                    <NButton
                      v-if="model.backgroundColor !== defaultCardBackground"
                      quaternary type="info" size="tiny"
                      @click="handleResetBackgroundColor"
                    >
                      {{ t('iconItem.resetBackground') }}
                    </NButton>
                  </div>
                </div>
              </NFormItem>

              <NGrid cols="2" :x-gap="10" item-responsive>
                <NGridItem span="2 500:1">
                  <NFormItem path="itemIconGroupId" :label="t('iconItem.iconGroup')">
                    <NSelect v-model:value="model.itemIconGroupId" :options="itemIconGroupOptions" size="small" />
                  </NFormItem>
                </NGridItem>

                <NGridItem span="2 500:1">
                  <NFormItem :label="t('iconItem.cardType')">
                    <NSelect v-model:value="model.cardType" :options="cardTypeOptions" size="small" />
                  </NFormItem>
                </NGridItem>
              </NGrid>

              <NFormItem path="onlyName">
                <template #label>
                  <div class="flex items-center">
                    <span>{{ t('common.onlyName') }}</span>
                    <span class="ml-1 text-xs opacity-60 cursor-help" :title="t('iconItem.onlyNameTip')">?</span>
                  </div>
                </template>
                <NInput
                  v-model:value="model.onlyName" type="text" show-count :maxlength="20"
                  :placeholder="t('iconItem.onlyNameTip')"
                />
              </NFormItem>

              <!-- 效果预览 -->
              <div class="mb-[10px]">
                <div class="text-xs text-slate-400 mb-[5px]">
                  {{ t('iconItem.preview') }}
                </div>
                <div class="flex items-center">
                  <div class="w-[220px]">
                    <AppIcon
                      :item-info="model"
                      :icon-text-color="panelState.panelConfig.iconTextColor"
                      :icon-text-info-hide-description="panelState.panelConfig.iconTextInfoHideDescription || false"
                      :icon-text-icon-hide-title="panelState.panelConfig.iconTextIconHideTitle || false"
                      :style="0"
                    />
                  </div>
                  <div class="w-[90px] ml-[10px]">
                    <AppIcon
                      :item-info="model"
                      :icon-text-color="panelState.panelConfig.iconTextColor"
                      :icon-text-info-hide-description="!panelState.panelConfig.iconTextInfoHideDescription"
                      :icon-text-icon-hide-title="panelState.panelConfig.iconTextIconHideTitle || false"
                      :style="1"
                    />
                  </div>
                </div>
              </div>
            </NCollapseItem>
          </NCollapse>
        </div>
      </NForm>
    </div>

    <template #footer>
      <NButton type="success" :loading="submitLoading" style="float: right;" @click="handleValidateButtonClick">
        {{ $t('common.save') }}
      </NButton>
    </template>
  </NModal>
</template>

<style scoped>
.card-bg-color-picker {
  width: 220px;
  min-width: 220px;
}
</style>
