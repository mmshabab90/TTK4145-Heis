#pragma once

#include "elev.h"

typedef enum {
    init,
    idle,
    doorOpen,
    moving,
    emergency,
    stoppedBetweenFloors
} State;

void fsm_init(void);

void fsm_event_buttonPressed(int floor, elev_button_type_t btn);
void fsm_event_stopButtonPressed(void);
void fsm_event_stopButtonReleased(void);
void fsm_event_arrivedAtFloor(int newFloor);
void fsm_event_timerHasTimedOut(void);
