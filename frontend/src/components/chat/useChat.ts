import {ref} from "vue";
import {Message, Session} from "./types";
import {createSession as _createSession, sendMessage} from './api'

const session = ref<Session|null>(null)

export function useChat() {

    const createSession = async ()=>{
        session.value = await _createSession();
    }
    const chat = async (message:Message)=>{

        if(session.value){

            if (!session.value.messages){
                session.value.messages = []
            }
            message.role = "user"
            session.value.messages.push(message)




            session.value.messages.push({
                role: "assistant",
                session: session.value.id,
                id: message.id,
                content: "",
                reason_content: "",
                request_id: message.request_id,
                meta:{
                    isThinking:false,
                    isThinkContentColspan:false
                }
            } as Message)

            let assistantMessage = session.value.messages[session.value.messages.length-1]
            
            const stream = sendMessage(message)

            const reader = stream.getReader();



            while (true) {
                const {done,value}  = await reader.read();
                if (done) {
                    break;
                }
                if(value){
                    if(value.error){
                        assistantMessage.error = value.error
                        break;
                    }
                    if(value.reason_content){
                        if(!assistantMessage.meta.isThinking){
                           assistantMessage.meta.isThinking = true
                        }
                        assistantMessage.reason_content += value.reason_content
                    }
                    if(value.content){
                        if(assistantMessage.meta?.isThinking){
                            assistantMessage.meta.isThinking = false
                        }
                        assistantMessage.content += value.content
                    }
                }
            }


            console.log(assistantMessage)

        }

    }

    return {
        chat,
        session,
        createSession
    }
}