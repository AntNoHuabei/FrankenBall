<script setup lang="ts">
import {h, onMounted, ref} from "vue";
import {useChat} from "./useChat";
import {Message as AIMessage} from "./types";
import {CircleUserIcon, MessageSquareIcon} from 'lucide-vue-next'
import { BulbOutlined, DownOutlined, LoadingOutlined, UpOutlined, UserOutlined,InfoCircleOutlined ,RocketOutlined,SmileOutlined,WarningOutlined} from '@ant-design/icons-vue';
import {Sender, Bubble, BubbleList,Prompts, type BubbleProps, PromptsProps,Welcome} from "ant-design-x-vue";
import {Empty,Avatar,Button,Space,Typography } from "ant-design-vue";

import { createMarkdownExit } from 'markdown-exit'
import { codeToHtml } from 'shiki'

// factory helper
const md = createMarkdownExit({
  html:true,
  linkify:true,
  async highlight(code, lang) {
    return codeToHtml(code, {lang, theme: 'nord'});
  }
})
defineOptions({ name: 'AXBubbleWithThinkSetup' });

const  inputValue = ref<string>('')

const {createSession,session, chat} = useChat();

const sendMessage = async (message:string)=>{
  if(session.value && inputValue.value){
    inputValue.value = ''
    await chat({
      id: Math.random().toString(),
      request_id: Math.random().toString(),
      content: message,
      role: "user",
      session:session.value?.id,
    })
  }

}

const processThink = (content:string)=>{
  if(content){
    let splits = content.split("\n");
    for(let index in splits){
      let item = splits[index];
      item = "> " + item;
      splits[index] = item;
    }
    return splits.join("\n")
  }else{
    return "";
  }
}

const renderMarkdown: BubbleProps['messageRender'] = (content) => {
  return h(Typography, null, () => h('div', { innerHTML: md.render(content) }));
};

onMounted(async ()=>{
  await createSession();
})


const items: PromptsProps['items'] = [
  {
    key: '1',
    icon: h(BulbOutlined, { style: { color: '#FFD700' } }),
    label: 'Ignite Your Creativity',
    description: 'Got any sparks for a new project?',
  },
  {
    key: '2',
    icon: h(InfoCircleOutlined, { style: { color: '#1890FF' } }),
    label: 'Uncover Background Info',
    description: 'Help me understand the background of this topic.',
  },
  {
    key: '3',
    icon: h(RocketOutlined, { style: { color: '#722ED1' } }),
    label: 'Efficiency Boost Battle',
    description: 'How can I work faster and better?',
  },
  {
    key: '4',
    icon: h(SmileOutlined, { style: { color: '#52C41A' } }),
    label: 'Tell me a Joke',
    description: 'Why do not ants get sick? Because they have tiny ant-bodies!',
  },
  {
    key: '5',
    icon: h(WarningOutlined, { style: { color: '#FF4D4F' } }),
    label: 'Common Issue Solutions',
    description: 'How to solve common issues? Share some tips!',
  },
];
</script>

<template>
  <div class="chatRoot">

    <div class="bubbleList">
      <template v-if="!session || !session.messages || session.messages.length === 0">
        <Welcome
            :icon="h('img', { src: '/logo.png', style: { width: '100%', height: '100%' } })"
            title="你好,我是Remo"
            description="你的个人助手,帮你记住各种工作细节"
        />
        <Prompts
            title="✨ 给你点灵感"
            :items="items"
            wrap
            :styles="{
      item: {
        flex: 'none',
        width: 'calc(50% - 6px)',
      },
    }"
        />
      </template>
      <template v-else v-for="msg in session.messages">

        <Bubble
            v-if="msg.role ==='user'"
            placement="end"
        :content="msg.content"
        >
          <template #avatar>
            <Avatar
              :icon="h(CircleUserIcon)"
              :size="35"
            />
          </template>
        </Bubble>

        <template v-if="msg.role ==='assistant'">

          <Bubble
              v-if="!msg.reason_content"
              variant="borderless"
              style="margin-top: -24px;"
              :typing="true"
              :content="msg.content"
              :message-render="renderMarkdown"
          />

          <Bubble
              placement="start"
              v-if="msg.reason_content"
              :avatar="{ icon: h(CircleUserIcon) }"
              :styles="{ footer: { marginTop: 0 } }"
          >
            <template #message>
              <Space>
                <BulbOutlined />
                <span>{{ msg.meta.isThinking ? "思考中..." : "已深度思考" }}</span>
                <Button
                    type="text"
                    size="small"
                    style="background: transparent;"
                    :icon="msg.meta.isThinkContentColspan ? h(UpOutlined) : h(DownOutlined)"
                    @click="() => {
                    msg.meta.isThinkContentColspan = !msg.meta.isThinkContentColspan
            }"
                />
              </Space>
            </template>
            <template #footer>
              <Space direction="vertical">
                <Bubble
                    v-show="!msg.meta.isThinkContentColspan"
                    variant="borderless"
                    :typing="false"
                    :content="processThink(msg.reason_content)"
                    :message-render="renderMarkdown"
                    @typing-complete="() => {
            }"
                />
                <LoadingOutlined v-if="msg.meta.isThinking" />
                <Bubble
                    variant="borderless"
                    style="margin-top: -24px;"
                    :typing="true"
                    :content="msg.content"
                    :message-render="renderMarkdown"
                />
              </Space>
            </template>
          </Bubble>
        </template>



      </template>
    </div>

    <Sender
    v-model:value="inputValue"
    @submit="sendMessage"
    @clear="()=>{
      inputValue = ''
    }"
    />
  </div>
</template>

<style lang="less" scoped>
.chatRoot{
  background-color: #ffffff;
  //all:initial;
  width: 100%;
  height: 100%;
  padding:15px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;

  .bubbleList{
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 10px;
    overflow-y: auto;
  }
}
</style>