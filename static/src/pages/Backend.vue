<script setup>
import { computed, inject, onMounted, ref } from "vue";
import ProviderEditor from "@/Components/Backend/ProviderEditor.vue";

const title = inject("title");
const activeProvider = ref("openai");

const providers = [
  {
    id: "openai",
    name: "OpenAI",
    vendor: "OpenAI",
    description: "OpenAI 原生或 OpenAI 兼容的主模型接口。",
    defaultEndpoint: "https://api.openai.com/v1/chat/completions",
    apiType: "openai",
    defaultTextModel: "gpt-4.1",
    defaultSummaryModel: "gpt-4.1-mini",
    defaultImageModel: "gpt-image-1",
    defaultImageSize: "1024x1024",
    defaultEmbeddingModel: "text-embedding-3-small",
    image: true,
    modelHints: ["gpt-4.1", "gpt-4.1-mini", "gpt-4o", "gpt-4o-mini", "o4-mini", "o3", "gpt-image-1", "text-embedding-3-small"],
  },
  {
    id: "deepseek",
    name: "DeepSeek",
    vendor: "DeepSeek",
    description: "DeepSeek V3/R1 等 OpenAI 兼容接口。",
    defaultEndpoint: "https://api.deepseek.com/chat/completions",
    apiType: "openai",
    defaultTextModel: "deepseek-chat",
    defaultSummaryModel: "deepseek-chat",
    embedding: false,
    modelHints: ["deepseek-chat", "deepseek-reasoner"],
  },
  {
    id: "alibaba",
    name: "阿里百炼",
    vendor: "Alibaba Cloud",
    description: "DashScope 兼容模式和通义万相图像接口配置。",
    defaultEndpoint: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
    apiType: "openai",
    defaultTextModel: "qwen-plus",
    defaultSummaryModel: "qwen-turbo",
    defaultImageModel: "wanx2.1-t2i-turbo",
    defaultImageSize: "1024*1024",
    image: true,
    embedding: false,
    modelHints: ["qwen-max", "qwen-plus", "qwen-turbo", "qwen-long", "deepseek-v3", "deepseek-r1", "wanx2.1-t2i-turbo"],
  },
  {
    id: "anthropic",
    name: "Claude",
    vendor: "Anthropic",
    description: "Anthropic Messages API，支持 Claude 系列模型。",
    defaultEndpoint: "https://api.anthropic.com/v1/messages",
    apiType: "anthropic",
    defaultTextModel: "claude-3-5-sonnet-latest",
    defaultSummaryModel: "claude-3-5-haiku-latest",
    embedding: false,
    defaultTopP: 1,
    modelHints: ["claude-3-5-sonnet-latest", "claude-3-5-haiku-latest", "claude-3-7-sonnet-latest"],
  },
  {
    id: "gemini",
    name: "Gemini",
    vendor: "Google",
    description: "Google Gemini API，支持从 Google 模型列表自动同步。",
    defaultEndpoint: "https://generativelanguage.googleapis.com/v1beta/models",
    apiType: "gemini",
    defaultTextModel: "gemini-2.0-flash",
    defaultSummaryModel: "gemini-2.0-flash-lite",
    defaultEmbeddingModel: "text-embedding-004",
    modelHints: ["gemini-2.0-flash", "gemini-2.0-flash-lite", "gemini-1.5-pro", "text-embedding-004"],
  },
  {
    id: "azure",
    name: "Azure OpenAI",
    vendor: "Microsoft Azure",
    description: "Azure OpenAI deployment 直连配置，模型名通常对应 deployment。",
    defaultEndpoint: "https://{resource}.openai.azure.com/openai/deployments/{deployment}/chat/completions?api-version=2024-02-01",
    apiType: "azure",
    defaultTextModel: "gpt-4o",
    defaultSummaryModel: "gpt-4o-mini",
    defaultEmbeddingModel: "text-embedding-3-small",
    canFetchModels: false,
    note: "Azure 通常需要填写完整 deployment URL；模型字段建议填写 deployment 名，列表拉取不适用于所有 Azure 资源。",
    modelHints: ["gpt-4o", "gpt-4o-mini", "gpt-4.1", "text-embedding-3-small"],
  },
  {
    id: "moonshot",
    name: "Moonshot Kimi",
    vendor: "Moonshot AI",
    description: "Moonshot 开放 API 或 Kimi for Coding 会员通道。",
    defaultEndpoint: "https://api.moonshot.cn/v1/chat/completions",
    apiType: "openai",
    defaultTextModel: "kimi-for-coding",
    defaultSummaryModel: "kimi-for-coding",
    embedding: false,
    defaultTemperature: 0.3,
    defaultTopP: 0.95,
    note: "Kimi for Coding 会员 Key 请将 API 地址改为 https://api.kimi.com/coding/v1/chat/completions，并在自定义 Header 中填写 User-Agent: claude-code/1.0 和 X-Client-Name: claude-code。",
    modelHints: ["kimi-for-coding", "kimi-k2.5", "kimi-latest", "kimi-k2-0711-preview", "moonshot-v1-8k", "moonshot-v1-32k", "moonshot-v1-128k"],
  },
  {
    id: "otherapi",
    name: "自定义兼容接口",
    vendor: "Custom",
    description: "任意 OpenAI 兼容网关、代理或私有部署。请求时请求头请设置 X-Provider: otherapi",
    defaultEndpoint: "",
    apiType: "openai",
    defaultTextModel: "",
    defaultSummaryModel: "",
    defaultEmbeddingModel: "",
    image: true,
    modelHints: [],
  },
];

const activeProviderConfig = computed(() => providers.find((item) => item.id === activeProvider.value) || providers[0]);

onMounted(() => {
  title.value = "接口设置";
});
</script>

<template>
  <mdui-layout-main class="backend-page">
    <div class="page-content">
      <header class="page-header">
        <div>
          <h1>模型接口</h1>
          <p>配置大模型提供方的鉴权、地址和请求协议。默认模型仅在请求未传 model 时兜底，模型列表可从远端同步。</p>
        </div>
        <div class="summary-strip">
          <span>{{ providers.length }} 个提供方</span>
          <span>文本 / 总结 / 图片 / Embedding</span>
        </div>
      </header>

      <div class="workspace">
        <aside class="provider-nav">
          <button v-for="provider in providers" :key="provider.id" class="provider-tab"
            :class="{ active: provider.id === activeProvider }" @click="activeProvider = provider.id">
            <span class="provider-name">{{ provider.name }}</span>
            <span class="provider-vendor">{{ provider.vendor }}</span>
          </button>
        </aside>

        <main class="editor-shell">
          <ProviderEditor v-for="provider in providers" :key="provider.id" :provider="provider"
            :active="activeProviderConfig.id === provider.id"></ProviderEditor>
        </main>
      </div>
    </div>
  </mdui-layout-main>
</template>

<style scoped>
.backend-page {
  background: #f6f8fb;
  box-sizing: border-box;
  display: block;
  min-height: 100%;
  width: 100%;
}

.page-content {
  box-sizing: border-box;
  padding: 0 1.5rem 1.5rem;
}

.page-header {
  align-items: flex-end;
  display: flex;
  gap: 1rem;
  justify-content: space-between;
  margin: 0 auto 1rem;
  max-width: 1280px;
}

h1,
p {
  margin: 0;
}

h1 {
  color: #202124;
  font-size: 1.9rem;
  font-weight: 760;
}

.page-header p {
  color: #5f6368;
  font-size: 0.95rem;
  line-height: 1.55;
  margin-top: 0.35rem;
  max-width: 50rem;
}

.summary-strip {
  align-items: flex-end;
  color: #5f6368;
  display: flex;
  flex-direction: column;
  font-size: 0.86rem;
  gap: 0.25rem;
  min-width: max-content;
}

.summary-strip span:first-child {
  color: #202124;
  font-size: 1rem;
  font-weight: 700;
}

.workspace {
  align-items: stretch;
  display: grid;
  gap: 1rem;
  grid-template-columns: 17rem minmax(0, 1fr);
  margin: 0 auto;
  max-width: 1280px;
}

.provider-nav {
  background: #fff;
  border: 1px solid #dde3ea;
  border-radius: 8px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  padding: 0.5rem;
}

.provider-tab {
  background: transparent;
  border: 1px solid transparent;
  border-radius: 7px;
  color: #3c4043;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  font: inherit;
  gap: 0.2rem;
  padding: 0.75rem;
  text-align: left;
}

.provider-tab:hover {
  background: #f1f5f9;
}

.provider-tab.active {
  background: #e8f0fe;
  border-color: #c7d7fe;
  color: #174ea6;
}

.provider-name {
  font-size: 0.96rem;
  font-weight: 700;
}

.provider-vendor {
  color: #6b7280;
  font-size: 0.78rem;
}

.editor-shell {
  background: #fff;
  border: 1px solid #dde3ea;
  border-radius: 8px;
  box-sizing: border-box;
  min-width: 0;
  padding: 1.25rem;
}

@media (max-width: 900px) {
  .backend-page {
    padding: 0;
  }

  .page-content {
    padding: 0 0.75rem 1rem;
  }

  .page-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .summary-strip {
    align-items: flex-start;
  }

  .workspace {
    grid-template-columns: 1fr;
  }

  .provider-nav {
    flex-direction: row;
    overflow-x: auto;
  }

  .provider-tab {
    min-width: 9rem;
  }
}
</style>
