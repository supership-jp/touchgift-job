package codes

// context.Context に設定するキー: type (ログをtype別に見れるようにする)
const KeyType = "type"

const TypeDeliveryStart = "delivery_start"
const TypeDeliveryEnd = "delivery_end"
const TypeDeliveryOperation = "delivery_operation"
const TypeDeliveryControl = "delivery_control"

const DetailShortage = "shortage"
const DetailExpended = "expended"
