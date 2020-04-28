namespace proto {

    let messages : Map<string, any>

    // 注册消息
    export function Register(msg : any) {
        let name = meta(()=>msg);
        messages[name] = msg;
    }

    let extractor = new RegExp("return (.*);");

    function meta<T>(name : ()=>T) {

        let array = extractor.exec(name + "");
        if (array == null) {
            throw new Error("function not match");
        }

        return array[1];
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