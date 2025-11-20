// frontend/src/utils/permMap.js
// TeamSpeak 3 完整权限映射表 (ID -> Name)
// 包含 Server, Channel, Client, Group, File Transfer 等常用权限

export default {
    // ============================
    // 1. 全局实例与虚拟服务器管理 (Global & Virtual Server)
    // ============================
    1: "b_serverinstance_help_view",
    2: "b_serverinstance_version_view",
    3: "b_serverinstance_info_view",
    4: "b_serverinstance_virtualserver_list",
    5: "b_serverinstance_binding_list",
    6: "b_serverinstance_permission_list",
    7: "b_serverinstance_permission_find",
    8: "b_virtualserver_create",
    9: "b_virtualserver_delete",
    10: "b_virtualserver_start_any",
    11: "b_virtualserver_stop_any",
    12: "b_virtualserver_change_machine_id",
    13: "b_virtualserver_change_template",
    14: "b_serverquery_login",
    15: "b_serverinstance_textmessage_send",
    16: "b_serverinstance_log_view",
    17: "b_serverinstance_log_add",
    18: "b_serverinstance_stop",
    19: "b_serverinstance_modify_settings",
    20: "b_serverinstance_modify_querygroup",
    21: "b_serverinstance_modify_templates",
    22: "b_virtualserver_select",
    23: "b_virtualserver_info_view",
    24: "b_virtualserver_connectioninfo_view",
    25: "b_virtualserver_channel_list",
    26: "b_virtualserver_channel_search",
    27: "b_virtualserver_client_list",
    28: "b_virtualserver_client_search",
    29: "b_virtualserver_client_dblist",
    30: "b_virtualserver_client_dbsearch",
    31: "b_virtualserver_client_dbinfo",
    32: "b_virtualserver_permission_find",
    33: "b_virtualserver_custom_search",
    34: "b_virtualserver_modify_name",
    35: "b_virtualserver_modify_welcomemessage",
    36: "b_virtualserver_modify_maxclients",
    37: "b_virtualserver_modify_reserved_slots",
    38: "b_virtualserver_modify_password",
    39: "b_virtualserver_modify_default_servergroup",
    40: "b_virtualserver_modify_default_channelgroup",
    41: "b_virtualserver_modify_hostbutton",
    42: "b_virtualserver_modify_hostbanner",
    43: "b_virtualserver_modify_gfx_interval",
    44: "b_virtualserver_modify_gfx_url",
    45: "b_virtualserver_modify_priority_speaker_dimm_modificator",
    46: "b_virtualserver_modify_log_settings",
    47: "b_virtualserver_modify_min_client_version",
    48: "b_virtualserver_modify_needed_identity_security_level",
    49: "b_virtualserver_modify_temporary_passwords",
    50: "b_virtualserver_modify_icon_id",

    // ============================
    // 2. 频道管理 (Channel Management)
    // ============================
    // 查看与基础
    55: "b_channel_info_view",
    56: "b_channel_create_child",
    57: "b_channel_create_permanent",
    58: "b_channel_create_semi_permanent",
    59: "b_channel_create_temporary",
    60: "b_channel_create_private",
    61: "b_channel_create_with_topic",
    62: "b_channel_create_with_description",
    63: "b_channel_create_with_password",

    // 频道修改 (Modify)
    64: "b_channel_modify_name",
    65: "b_channel_modify_topic",
    66: "b_channel_modify_description",
    67: "b_channel_modify_password",
    68: "b_channel_modify_codec",
    69: "b_channel_modify_codec_quality",
    70: "b_channel_modify_maxclients",
    71: "b_channel_modify_maxfamilyclients",
    72: "b_channel_modify_sortorder",
    73: "b_channel_modify_default",
    74: "b_channel_modify_needed_talk_power",
    75: "b_channel_modify_make_permanent",
    76: "b_channel_modify_make_semi_permanent",
    77: "b_channel_modify_make_temporary",

    // 频道权限 (Power)
    78: "i_channel_create_modify_with_codec_maxquality",
    79: "i_channel_create_modify_with_codec_latency_factor_min",
    80: "i_channel_modify_power",
    81: "i_channel_needed_modify_power",
    82: "b_channel_modify_parent",
    83: "b_channel_delete_permanent",
    84: "b_channel_delete_semi_permanent",
    85: "b_channel_delete_temporary",
    86: "b_channel_delete_flag_force",
    87: "i_channel_delete_power",
    88: "i_channel_needed_delete_power",

    // 频道访问 (Access)
    89: "b_channel_join_permanent",
    90: "b_channel_join_semi_permanent",
    91: "b_channel_join_temporary",
    92: "b_channel_join_ignore_password",
    93: "b_channel_join_ignore_maxclients",
    94: "i_channel_join_power",
    95: "i_channel_needed_join_power",
    96: "i_channel_subscribe_power",
    97: "i_channel_needed_subscribe_power",
    98: "i_channel_description_view_power",
    99: "i_channel_needed_description_view_power",

    // ============================
    // 3. 客户端管理与通信 (Client & Communication)
    // ============================
    // 踢出与封禁 (Kick & Ban)
    100: "i_client_kick_from_server_power",
    101: "i_client_needed_kick_from_server_power",
    102: "i_client_kick_from_channel_power",
    103: "i_client_needed_kick_from_channel_power",
    104: "i_client_ban_power",
    105: "i_client_needed_ban_power",
    106: "i_client_ban_max_bantime", // -1 = 永久
    107: "b_client_ban_list",
    108: "b_client_ban_create",
    109: "b_client_ban_delete",
    110: "b_client_ignore_ban",

    // 移动与基本交互 (Move & Talk)
    111: "i_client_move_power",
    112: "i_client_needed_move_power",
    113: "i_client_talk_power",
    114: "i_client_needed_talk_power",
    115: "b_client_is_priority_speaker",
    116: "b_client_force_push_to_talk",
    117: "i_client_poke_power",
    118: "i_client_needed_poke_power",
    119: "i_client_whisper_power",
    120: "i_client_needed_whisper_power",
    121: "b_client_ignore_whisper",

    // 消息与投诉 (Message & Complaint)
    122: "i_client_private_text_message_power",
    123: "i_client_needed_private_text_message_power",
    124: "b_client_server_text_message_send",
    125: "b_client_channel_text_message_send",
    126: "b_client_offline_text_message_send",
    127: "i_client_complaint_power",
    128: "i_client_needed_complaint_power",
    129: "b_client_complaint_list",
    130: "b_client_complaint_delete_own",
    131: "b_client_complaint_delete",

    // ============================
    // 4. 组管理 (Group Management)
    // ============================
    132: "b_virtualserver_servergroup_list",
    133: "b_virtualserver_servergroup_permission_list",
    134: "b_virtualserver_servergroup_client_list",
    135: "i_server_group_modify_power",
    136: "i_server_group_needed_modify_power",
    137: "i_server_group_member_add_power",
    138: "i_server_group_needed_member_add_power",
    139: "i_server_group_member_remove_power",
    140: "i_server_group_needed_member_remove_power",
    141: "b_server_group_create", // 非标准 ID，视版本而定
    142: "b_server_group_delete",
    143: "i_group_sort_id",
    144: "i_group_show_name_in_tree",
    145: "i_group_auto_update_type",
    146: "i_group_auto_update_max_value",
    147: "b_group_is_permanent",

    // ============================
    // 5. 文件传输 (File Transfer)
    // ============================
    150: "b_ft_ignore_password",
    151: "b_ft_transfer_list",
    152: "i_ft_file_upload_power",
    153: "i_ft_needed_file_upload_power",
    154: "i_ft_file_download_power",
    155: "i_ft_needed_file_download_power",
    156: "i_ft_file_delete_power",
    157: "i_ft_needed_file_delete_power",
    158: "i_ft_file_rename_power",
    159: "i_ft_needed_file_rename_power",
    160: "i_ft_file_browse_power",
    161: "i_ft_needed_file_browse_power",
    162: "i_ft_directory_create_power",
    163: "i_ft_needed_directory_create_power",
    164: "i_ft_quota_mb_download_per_client",
    165: "i_ft_quota_mb_upload_per_client",

    // ============================
    // 6. 杂项与图标 (Misc & Icon)
    // ============================
    170: "i_icon_id",             // 组图标ID
    171: "b_icon_manage",         // 上传/删除图标
    172: "b_client_info_view",    // 查看连接信息
    173: "b_client_permission_overview_view", // 查看权限概览
    174: "b_client_remoteaddress_view",       // 查看IP地址
    175: "b_client_custom_info_view",
    176: "i_client_max_avatar_filesize",
    177: "i_client_max_channel_subscriptions",

    // ============================
    // 7. 常用缺失 ID 补充
    // ============================
    // 下面这些ID可能因 TS3 版本差异而不同，但通常位于高位区间
    // 请根据你实际看到的“灰色ID”在列表中核对，TS3 ID 不是完全连续的
    65534: "i_needed_modify_power_client_config_delete",
    65535: "b_client_skip_channelgroup_permissions"
}