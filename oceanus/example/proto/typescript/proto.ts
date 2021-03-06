namespace proto {

    let messages : Map<string, any>

    export interface Message {
        Meta() : string;
    }

    // 注册消息
    export function Register(name : string, type : any) {

    }

    const separator : string = "-"

    // 反序列化消息
    export function Unmarshal(raw : string) : any {

        let bp = raw.split(separator, 2)
        if (bp.length != 2) {
            return null;
        }

        let type = messages[bp[0]]
        if (type == undefined) {
            return null;
        }

        let object = JSON.parse(bp[1])
        return Object.assign(new type, object)
    }
}