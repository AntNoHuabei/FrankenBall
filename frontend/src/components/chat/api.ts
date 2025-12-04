import {Message, Session} from "./types";
import {fetchEventSource} from "@microsoft/fetch-event-source";


export const createSession = ():Promise<Session>=>{

    return new Promise((resolve,reject)=>{

        const result = fetch("/api/session/create",{
            method: 'POST',
        })
        result.then(res=>{
            res.json().then(data=>{
                console.log(data)
                if(data.code === 200 && data.data){
                    resolve(data.data as Session);
                }else{
                    reject(data.message);
                }
            })
        }).catch( err=>{
            reject(err);
        })
    })

}

export const sendMessage = (message:Message):ReadableStream<Message> => {

    return new ReadableStream({
        start(controller) {

            fetchEventSource('http://localhost:9980/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                openWhenHidden: false,
                body: JSON.stringify({
                    message: message.content,
                    session: message.session
                }),
                onmessage(event) {
                    try {
                        const message = JSON.parse(event.data) as Message
                            controller.enqueue(message)
                    }catch ( e){
                        console.log( e)
                    }
                },
            }).then(r  =>{
                controller.close();
            }).catch( err=>{
                controller.enqueue({
                    error:err.message || err || "未知错误"
                } as Message)
            });


        },
    });
}