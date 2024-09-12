package codes

// context.Context に設定するキー: type (ログをtype別に見れるようにする)
const KeyType = "type"

const TypeDeliveryStart = "delivery_start"
const TypeDeliveryEnd = "delivery_end"
const TypeDeliveryOperation = "delivery_operation"
const TypeDeliveryControl = "delivery_control"

const DetailShortage = "shortage"
const DetailExpended = "expended"

const StatusStart = "start"
const StatusStarted = "started"
const StatusWarmup = "warmup"
const StatusResume = "resume"
const StatusPause = "pause"
const StatusPaused = "paused"
const StatusStopped = "stopped"
const StatusEnd = "end"
const StatusEnded = "ended"
const StatusSuspend = "suspend"
const StatusConfigured = "configured"
const StatusStop = "stop"
const StatusTerminate = "terminate"
