export interface MessageMeta{
    isThinking?:boolean
    isThinkContentColspan?:boolean
}
export interface Message{
    id: string;
    content: string;
    reason_content?:string ;
    role: 'system' | 'user' | 'assistant';
    request_id: string;
    session:string;
    error?:string;
    meta?:MessageMeta
}


export interface Session{
    id: string;
    title:string;
    messages: Message[];
}