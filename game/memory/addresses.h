#ifndef ADDRESSES_H
#define ADDRESSES_H

#include <stdint.h>

const int SIGNAL_EVENT_FUN_PTR = 0x703e50;
const int CLICK_TO_MOVE_FIX_PTR = 0x860A90;
const uint32_t UNIT_NAME_OFFSET = 0xB30;
const int* NAME_LIST_POINTER = (int*) 0xC0E230;
const int NAME_PTR_GUID_OFFSET = 0xC;
const int NAME_PTR_NAME_OFFSET = 0x14;
const int DESCRIPTOR_OFFSET = 0x8;
const int CURRENT_HEALTH_OFFSET = 0x58;
const int MAX_HEALTH_OFFSET = 0x70;
const int CURRENT_MANA_OFFSET = 0x5C;
const int MAX_MANA_OFFSET = 0x74;
const int TARGET_GUID_OFFSET = 0x40;
const int RAGE_OFFSET = 0x60;
const int ENERGY_OFFSET = 0x68;
const int FACTION_ID_OFFSET = 0x8C;
const int UNIT_FLAGS_OFFSET = 0xB8;
const int BUFFS_BASE_OFFSET = 0xBC;
const int DEBUFFS_BASE_OFFSET = 0x13C;
const int LEVEL_OFFSET = 0x88;
const int DYNAMIC_FLAGS_OFFSET = 0x23C;

const int CURRENT_SPELLCAST_OFFSET = 0xC8C;
const int POS_X_OFFSET = 0x9B8;
const int POS_Y_OFFSET = 0x9BC;
const int POS_Z_OFFSET = 0x9C0;

#endif