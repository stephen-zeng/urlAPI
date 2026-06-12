<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import Cookies from "js-cookie";
import { Post } from "@/js/fetch.js";
import { Notification, Setting } from "@/js/util.js";

const props = defineProps({
  provider: {
    type: Object,
    required: true,
  },
  active: {
    type: Boolean,
    default: false,
  },
});

const settings = ref(null);
const loading = ref(false);
const saving = ref(false);
const modelLoading = ref(false);
const modelError = ref("");
const fetchedModels = ref([]);
const customHeaderText = ref("");

const apiTypeOptions = [
  { value: "openai", label: "OpenAI 兼容" },
  { value: "anthropic", label: "Anthropic Messages" },
  { value: "gemini", label: "Google Gemini" },
  { value: "azure", label: "Azure OpenAI" },
];

const form = reactive({
  api_key: "",
  endpoint: "",
  api_type: "",
  text_model: "",
  summary_model: "",
  image_model: "",
  image_size: "",
  embedding_model: "",
  temperature: 1,
  top_p: 1,
  max_tokens: 0,
  presence_penalty: 0,
  frequency_penalty: 0,
  enabled: true,
});

const headerSummary = computed(() => {
  if (!settings.value) return "未加载";
  if (form.enabled === false) return "已停用";
  if (!settings.value.api_key_set && !form.api_key) return "缺少 Key";
  return form.text_model || "未选择模型";
});

const modelOptions = computed(() => {
  const pool = [
    ...fetchedModels.value,
    ...(props.provider.modelHints || []),
    form.text_model,
    form.summary_model,
  ];
  return [...new Set(pool.map((item) => String(item || "").trim()).filter(Boolean))];
});

const canFetchModels = computed(() => props.provider.canFetchModels !== false);

onMounted(() => {
  loadSettings();
});

async function loadSettings() {
  loading.value = true;
  try {
    const data = await Setting("fetchSettings", props.provider.id);
    if (!data) return;
    const textModel = data.text_model || data.summary_model || props.provider.defaultTextModel || "";
    const apiType = normalizeAPIType(data.api_type || props.provider.apiType || "openai");
    settings.value = data;
    Object.assign(form, {
      api_key: "",
      endpoint: data.endpoint || props.provider.defaultEndpoint || "",
      api_type: apiType,
      text_model: textModel,
      summary_model: textModel,
      image_model: data.image_model || props.provider.defaultImageModel || "",
      image_size: data.image_size || props.provider.defaultImageSize || "",
      embedding_model: data.embedding_model || props.provider.defaultEmbeddingModel || "",
      temperature: numberOr(data.temperature, props.provider.defaultTemperature ?? 1),
      top_p: numberOr(data.top_p, props.provider.defaultTopP ?? 1),
      max_tokens: Number(data.max_tokens || 0),
      presence_penalty: numberOr(data.presence_penalty, 0),
      frequency_penalty: numberOr(data.frequency_penalty, 0),
      enabled: data.enabled !== false,
    });
    customHeaderText.value = formatHeaders(data.custom_headers);
  } finally {
    loading.value = false;
  }
}

async function saveSettings() {
  saving.value = true;
  try {
    const headers = parseHeaders(customHeaderText.value);
    if (headers === null) return;
    const payload = {
      ...form,
      summary_model: form.text_model,
      temperature: Number(form.temperature),
      top_p: Number(form.top_p),
      max_tokens: Number(form.max_tokens || 0),
      presence_penalty: Number(form.presence_penalty),
      frequency_penalty: Number(form.frequency_penalty),
      custom_headers: headers,
    };
    if (!payload.api_key) {
      delete payload.api_key;
    }
    await Setting("editSettings", props.provider.id, payload);
    await loadSettings();
  } finally {
    saving.value = false;
  }
}

async function fetchModels() {
  modelLoading.value = true;
  modelError.value = "";
  try {
    const resp = await Post({
      Token: Cookies.get("token"),
      Send: {
        operation: "fetchProviderModels",
        setting_part: `provider.${props.provider.id}`,
      },
    });
    if (resp?.error) {
      modelError.value = resp.error;
      Notification(resp.error);
      return;
    }
    fetchedModels.value = resp?.setting_body?.models || resp?.models || [];
    if (fetchedModels.value.length === 0) {
      modelError.value = "远端没有返回模型列表，请检查 API Key、地址和接口类型";
    }
  } catch (err) {
    modelError.value = err.message;
    Notification("获取模型列表失败: " + err.message);
  } finally {
    modelLoading.value = false;
  }
}

function setModel(value) {
  form.text_model = value;
  form.summary_model = value;
}

function applyPreset(preset) {
  form.temperature = preset.temperature;
  form.top_p = preset.top_p;
  form.presence_penalty = preset.presence_penalty;
  form.frequency_penalty = preset.frequency_penalty;
}

function numberOr(value, fallback) {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function normalizeAPIType(value) {
  return apiTypeOptions.some((option) => option.value === value) ? value : "openai";
}

function formatHeaders(headers) {
  if (!headers || Object.keys(headers).length === 0) return "";
  return Object.entries(headers).map(([key, value]) => `${key}: ${value}`).join("\n");
}

function parseHeaders(text) {
  const headers = {};
  for (const line of text.split("\n")) {
    const trimmed = line.trim();
    if (!trimmed) continue;
    const index = trimmed.indexOf(":");
    if (index <= 0) {
      Notification("自定义 Header 格式应为 Key: Value");
      return null;
    }
    headers[trimmed.slice(0, index).trim()] = trimmed.slice(index + 1).trim();
  }
  return headers;
}
</script>

<template>
  <section v-show="active" class="provider-editor">
    <div class="editor-header">
      <div>
        <div class="eyebrow">{{ provider.vendor }}</div>
        <h2>{{ provider.name }}</h2>
        <p>{{ provider.description }}</p>
      </div>
      <div class="header-actions">
        <span class="status-pill" :class="{ off: !form.enabled }">{{ headerSummary }}</span>
        <mdui-switch :checked="form.enabled" @change="form.enabled = $event.target.checked"></mdui-switch>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <mdui-circular-progress></mdui-circular-progress>
      <span>加载配置中</span>
    </div>

    <div v-else class="editor-grid">
      <section class="panel credentials-panel">
        <div class="panel-title">
          <mdui-icon name="key"></mdui-icon>
          <div>
            <h3>连接</h3>
            <p>{{ settings?.api_key_set ? "已保存 API Key，留空不会覆盖" : "尚未保存 API Key" }}；API 类型表示请求协议</p>
          </div>
        </div>
        <mdui-text-field class="field" variant="outlined" label="API Key" type="password" toggle-password
          :placeholder="settings?.api_key_set ? '已保存，输入新值才会替换' : '输入 API Key'"
          :value="form.api_key" @input="form.api_key = $event.target.value"></mdui-text-field>
        <mdui-text-field class="field" variant="outlined" label="API 地址" :value="form.endpoint"
          @input="form.endpoint = $event.target.value"></mdui-text-field>
        <mdui-select class="field" variant="outlined" label="API 类型" :value="form.api_type"
          @change="form.api_type = $event.target.value">
          <mdui-menu-item v-for="option in apiTypeOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </mdui-menu-item>
        </mdui-select>
        <p v-if="provider.note" class="note">{{ provider.note }}</p>
      </section>

      <section class="panel model-panel">
        <div class="panel-title model-title">
          <mdui-icon name="hub"></mdui-icon>
          <div>
            <h3>模型</h3>
            <p>请求体中的 model 优先生效；这里的值仅作为请求未指定 model 时的兜底</p>
          </div>
          <mdui-button class="inline-button" variant="tonal" icon="sync" :loading="modelLoading"
            :disabled="!canFetchModels" @click="fetchModels">拉取模型</mdui-button>
        </div>
        <p v-if="modelError" class="note warning">{{ modelError }}</p>
        <div class="model-chips" v-if="modelOptions.length">
          <mdui-chip v-for="model in modelOptions.slice(0, 18)" :key="model" :selected="model === form.text_model"
            @click="setModel(model)">
            {{ model }}
          </mdui-chip>
        </div>
        <mdui-select v-if="modelOptions.length" class="field" variant="outlined" label="默认模型（可选）"
          :value="form.text_model" @change="setModel($event.target.value)">
          <mdui-menu-item v-for="model in modelOptions" :key="model" :value="model">{{ model }}</mdui-menu-item>
        </mdui-select>
        <mdui-text-field v-else class="field" variant="outlined" label="默认模型（可选）"
          :value="form.text_model" @input="setModel($event.target.value)"></mdui-text-field>
      </section>

      <section class="panel parameters-panel">
        <div class="panel-title">
          <mdui-icon name="tune"></mdui-icon>
          <div>
            <h3>参数控制</h3>
            <p>仅用于站内生成任务；OpenAI 兼容 API 请求中的同名参数优先生效</p>
          </div>
        </div>
        <div class="preset-row">
          <mdui-button class="inline-button" variant="outlined" @click="applyPreset({ temperature: 0.2, top_p: 0.8, presence_penalty: 0, frequency_penalty: 0 })">稳定</mdui-button>
          <mdui-button class="inline-button" variant="outlined" @click="applyPreset({ temperature: 0.7, top_p: 0.95, presence_penalty: 0, frequency_penalty: 0 })">均衡</mdui-button>
          <mdui-button class="inline-button" variant="outlined" @click="applyPreset({ temperature: 1, top_p: 1, presence_penalty: 0.2, frequency_penalty: 0.2 })">发散</mdui-button>
        </div>
        <div class="slider-field">
          <label>Temperature <strong>{{ Number(form.temperature).toFixed(2) }}</strong></label>
          <mdui-slider :value="form.temperature" min="0" max="2" step="0.05"
            @input="form.temperature = Number($event.target.value)"></mdui-slider>
        </div>
        <div class="slider-field">
          <label>Top P <strong>{{ Number(form.top_p).toFixed(2) }}</strong></label>
          <mdui-slider :value="form.top_p" min="0" max="1" step="0.01"
            @input="form.top_p = Number($event.target.value)"></mdui-slider>
        </div>
        <div class="numeric-grid">
          <mdui-text-field class="field" variant="outlined" type="number" label="Max Tokens，0 为不限"
            :value="form.max_tokens" @input="form.max_tokens = Number($event.target.value || 0)"></mdui-text-field>
          <mdui-text-field class="field" variant="outlined" type="number" label="Presence Penalty"
            :value="form.presence_penalty" min="-2" max="2" step="0.1"
            @input="form.presence_penalty = Number($event.target.value || 0)"></mdui-text-field>
          <mdui-text-field class="field" variant="outlined" type="number" label="Frequency Penalty"
            :value="form.frequency_penalty" min="-2" max="2" step="0.1"
            @input="form.frequency_penalty = Number($event.target.value || 0)"></mdui-text-field>
        </div>
      </section>

      <section class="panel headers-panel">
        <div class="panel-title">
          <mdui-icon name="notes"></mdui-icon>
          <div>
            <h3>自定义 Header</h3>
            <p>每行一个 Key: Value；可用于代理、组织 ID，或覆盖 User-Agent 等上游要求的客户端 Header</p>
          </div>
        </div>
        <mdui-text-field class="field textarea" variant="outlined" rows="5" autosize label="Headers"
          :value="customHeaderText" @input="customHeaderText = $event.target.value"></mdui-text-field>
      </section>
    </div>

    <div class="save-bar">
      <mdui-button class="inline-button" variant="text" icon="refresh" @click="loadSettings">重载</mdui-button>
      <mdui-button class="save-button" variant="filled" icon="save" :loading="saving" @click="saveSettings">保存配置</mdui-button>
    </div>
  </section>
</template>

<style scoped>
.provider-editor {
  min-width: 0;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.eyebrow {
  color: #5f6368;
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0;
  margin-bottom: 0.25rem;
}

h2,
h3,
p {
  margin: 0;
}

h2 {
  color: #202124;
  font-size: 1.65rem;
  font-weight: 720;
}

.editor-header p,
.panel-title p,
.note {
  color: #5f6368;
  font-size: 0.88rem;
  line-height: 1.5;
}

.header-actions {
  align-items: center;
  display: flex;
  gap: 0.75rem;
  min-width: max-content;
}

.status-pill {
  background: #e8f0fe;
  border: 1px solid #c7d7fe;
  border-radius: 999px;
  color: #174ea6;
  display: inline-flex;
  font-size: 0.82rem;
  font-weight: 650;
  padding: 0.35rem 0.7rem;
}

.status-pill.off {
  background: #fce8e6;
  border-color: #fad2cf;
  color: #a50e0e;
}

.loading-state {
  align-items: center;
  color: #5f6368;
  display: flex;
  gap: 0.75rem;
  min-height: 12rem;
}

.editor-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: minmax(17rem, 0.9fr) minmax(24rem, 1.4fr);
}

.panel {
  background: #fff;
  border: 1px solid #dde3ea;
  border-radius: 8px;
  box-sizing: border-box;
  padding: 1rem;
}

.model-panel,
.parameters-panel {
  min-width: 0;
}

.panel-title {
  align-items: flex-start;
  display: flex;
  gap: 0.7rem;
  margin-bottom: 1rem;
}

.panel-title mdui-icon {
  color: #1a73e8;
  margin-top: 0.1rem;
}

.model-title {
  align-items: center;
}

.model-title > div {
  flex: 1;
}

.field {
  box-sizing: border-box;
  margin: 0 0 0.75rem;
  width: 100%;
}

.note {
  background: #f8fafc;
  border-left: 3px solid #9aa0a6;
  margin-top: 0.25rem;
  padding: 0.55rem 0.7rem;
}

.note.warning {
  background: #fef7e0;
  border-color: #f9ab00;
  color: #8a5a00;
  margin-bottom: 0.75rem;
}

.model-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-bottom: 0.8rem;
  max-height: 7.5rem;
  overflow: auto;
}

.numeric-grid {
  display: grid;
  gap: 0.75rem;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.slider-field {
  margin-bottom: 0.9rem;
}

.slider-field label {
  color: #3c4043;
  display: flex;
  font-size: 0.88rem;
  justify-content: space-between;
  margin-bottom: 0.2rem;
}

.preset-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.inline-button {
  width: auto;
}

.textarea {
  margin-bottom: 0;
}

.save-bar {
  align-items: center;
  background: #f8fafc;
  border: 1px solid #dde3ea;
  border-radius: 8px;
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 1rem;
  padding: 0.75rem;
}

.save-button {
  width: auto;
}

@media (max-width: 980px) {
  .editor-header,
  .header-actions {
    align-items: flex-start;
    flex-direction: column;
  }

  .editor-grid,
  .numeric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
