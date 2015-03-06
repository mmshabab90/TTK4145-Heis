#pragma once

#include "elev.h"   // for types only


void orders_add(int floor, elev_button_type_t btn);

elev_motor_direction_t orders_chooseDirn(int currFloor, elev_motor_direction_t currDirn);

int orders_shouldStop(int floor, elev_motor_direction_t dirn);

void orders_clearOrdersAt(int floor);

void orders_removeAll(void);

int orders_read(int floor, elev_button_type_t btn);
