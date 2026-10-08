<script setup lang="ts">
import { useDesign } from '@/hooks/web/useDesign'
import { ElButton, ComponentSize, ButtonType } from 'element-plus'
import { PropType, Component } from 'vue'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('button')

const props = defineProps({
  size: {
    type: String as PropType<ComponentSize>,
    default: undefined
  },
  type: {
    type: String as PropType<ButtonType>,
    default: 'default'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  plain: {
    type: Boolean,
    default: false
  },
  text: {
    type: Boolean,
    default: false
  },
  bg: {
    type: Boolean,
    default: false
  },
  link: {
    type: Boolean,
    default: false
  },
  round: {
    type: Boolean,
    default: false
  },
  circle: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  loadingIcon: {
    type: [String, Object] as PropType<String | Component>,
    default: undefined
  },
  icon: {
    type: [String, Object] as PropType<String | Component>,
    default: undefined
  },
  autofocus: {
    type: Boolean,
    default: false
  },
  nativeType: {
    type: String as PropType<'button' | 'submit' | 'reset'>,
    default: 'button'
  },
  autoInsertSpace: {
    type: Boolean,
    default: false
  },
  color: {
    type: String,
    default: ''
  },
  darker: {
    type: Boolean,
    default: false
  },
  tag: {
    type: [String, Object] as PropType<String | Component>,
    default: 'button'
  }
})

const emits = defineEmits(['click'])

</script>

<template>
  <ElButton
    :class="prefixCls"
    v-bind="{ ...props }"
    @click="() => emits('click')"
  >
    <slot></slot>
    <slot name="icon"></slot>
    <slot name="loading"></slot>
  </ElButton>
</template>
