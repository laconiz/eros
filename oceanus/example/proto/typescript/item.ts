namespace proto {

    // 物品ID
    export type ItemID = string;

    // 物品定义
    export class Item {
        ID : ItemID; // ID
        Num : number; // 数量
    }

    // 金币
    export const ItemCoin : ItemID = "coin";
    // 券
    export const ItemTicket : ItemID = "ticket";

    // 物品列表
    export type Items = Map<ItemID, number>;

    // 物品列表更新消息
    export class ItemsACK {
        Items : Items;
    }

    Register("ItemsACK", ItemsACK)
}