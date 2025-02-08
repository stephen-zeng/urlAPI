<script setup>
import {inject, onMounted, ref} from "vue";

const url = inject("url");
const type = inject("type");
const show = inject("show");

const host = ref("https://");
const operation = ref("");
const config = ref("");
const configs = ref("");

onMounted(()=>{
  host.value = 'https://' + document.location.host;
})

function update(id, value) {
  show.value = false;
  if (id == "type") type.value = value;
  if (id == "host") host.value = value;
  if (id == "operation") {
    operation.value = value;
    configs.value = '';
  }
  if (id == "configs") {
    config.value='';
    if (configs.value == "") configs.value += value;
    else configs.value += '&' + value;
  }
  url.value = host.value + '/' + operation.value + '?' + configs.value;
}

</script>

<template>
  <mdui-card variant="outlined" class="workshop-control-card">
    <h1>工作台</h1>
    <p>测试API的使用情况，请将API的域名添加到Referer白名单中</p>
    <mdui-radio-group :value="type" @change="update('type',$event.target.value)">
      <mdui-radio value="img">图片 (img)</mdui-radio>
      <mdui-radio value="txt">文字 (iframe)</mdui-radio>
    </mdui-radio-group>
    <mdui-text-field variant="outlined" label="Host"
                     :value="host" @input="update('host', $event.target.value)"></mdui-text-field>
    <mdui-text-field variant="outlined" label="Operation"
                     :value="operation" @input="update('operation', $event.target.value)"></mdui-text-field>
    <mdui-text-field variant="outlined" label="Config (xx=xx)"
                     :value="config" @input="config = $event.target.value"
                     @keydown.enter="update('configs', config)">
      <mdui-button-icon slot="end-icon" icon="add"
                        @click="update('configs', config)"></mdui-button-icon>
    </mdui-text-field>
    <mdui-text-field variant="outlined" label="URL" autosize
                     readonly :value="url"></mdui-text-field>
    <mdui-button @click="show=true">Go</mdui-button>
  </mdui-card>
</template>

<style scoped>
</style>