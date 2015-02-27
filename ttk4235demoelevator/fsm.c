#include <stdio.h>

#include "fsm.h"
#include "orders.h"
#include "timer.h"
#include "elev.h"



static State                    state;
static elev_motor_direction_t   dirn;
static int                      floor;
static elev_motor_direction_t   departDirn;


const double                    doorOpenTime    = 3.0;

char* State_toString(State s){
    switch(s){
    case init:                      return "init";
    case idle:                      return "idle";
    case doorOpen:                  return "doorOpen";
    case moving:                    return "moving";
    case emergency:                 return "emergency";
    case stoppedBetweenFloors:      return "stoppedBetweenFloors";
    default:                        return "State_WAT";
    }
}

char* Button_toString(elev_button_type_t btn){
    switch(btn){
    case BUTTON_CALL_UP:    return "UP";
    case BUTTON_CALL_DOWN:  return "DOWN";
    case BUTTON_COMMAND:    return "COMMAND";
    default:                return "Button_WAT";
    }
}

void setAllLights(){
    for(int floor = 0; floor < N_FLOORS; floor++){
        for(int btn = 0; btn < 3; btn++){
            if( (btn == BUTTON_CALL_UP      &&  floor == N_FLOORS-1) ||
                (btn == BUTTON_CALL_DOWN    &&  floor == 0))
            {
                continue;
            }
            elev_set_button_lamp(btn, floor, orders_read(floor, btn));
        }
    }
}


void fsm_init(void){
    state       = init;
    dirn        = DIRN_STOP;
    floor       = -1;
    departDirn  = DIRN_DOWN;
    orders_removeAll();
}






void fsm_event_buttonPressed(int btn_floor, elev_button_type_t btn_type){
    printf("%s(%d, %s)\n  State:     %-15.15s |  Floor: %d  |  Dirn: %d\n",  
        __FUNCTION__, btn_floor, Button_toString(btn_type), 
        State_toString(state), floor, dirn
    );
    
    switch(state){
    case init:
    case emergency:
        break;
        
    case doorOpen:
        if(floor == btn_floor){
            timer_start(doorOpenTime);
        } else {
            orders_add(btn_floor, btn_type);
        }
        break;
    case moving:
        orders_add(btn_floor, btn_type);
        break;
        
    case idle:
        orders_add(btn_floor, btn_type);
        
        dirn = orders_chooseDirn(floor, dirn);
        
        if(dirn == DIRN_STOP){
            elev_set_door_open_lamp(1);
            orders_clearOrdersAt(floor);
            timer_start(doorOpenTime);
            state = doorOpen;
        } else {
            elev_set_motor_direction(dirn);
            state = moving;
            departDirn = dirn;
        }
    
        break;
        
    case stoppedBetweenFloors:
        orders_add(btn_floor, btn_type);
        
        dirn = orders_chooseDirn(floor, departDirn);
        elev_set_motor_direction(dirn);
    
        state = moving;
    
        break;
        
    default:
        printf("  \7WAT at %s:%d: State \"%s\" is invalid!\n",
            __FUNCTION__, __LINE__,
            State_toString(state)
        );
        
    }
    
    setAllLights();
    
    printf("  New state: %s\n", State_toString(state));
}




void fsm_event_stopButtonPressed(void){
    printf("%s\n  State:     %-15.15s |  Floor: %d  |  Dirn: %d\n",  
        __FUNCTION__, 
        State_toString(state), floor, dirn
    );
    
    switch(state){
    case init:
    case emergency:
        break;
        
    case idle:
    case doorOpen:
    case moving:
    case stoppedBetweenFloors:
        elev_set_motor_direction(DIRN_STOP);
        elev_set_stop_lamp(1);
        orders_removeAll();
        timer_stop();
        
        setAllLights();
        
        if(elev_get_floor_sensor_signal() != -1){
            elev_set_door_open_lamp(1);
        }
        
        state = emergency;
        
        break;
        
    default:
        printf("  \7WAT at %s:%d: State \"%s\" is invalid!\n",
            __FUNCTION__, __LINE__,
            State_toString(state)
        );        
    }
    
    printf("  New state: %s\n", State_toString(state));
}




void fsm_event_stopButtonReleased(void){
    printf("%s\n  State:     %-15.15s |  Floor: %d  |  Dirn: %d\n",  
        __FUNCTION__, 
        State_toString(state), floor, dirn
    );
    
    switch(state){
    case emergency:
        elev_set_stop_lamp(0);
        
        if(elev_get_floor_sensor_signal() == -1){
            state = stoppedBetweenFloors;
        } else {
            dirn = DIRN_STOP;
            timer_start(doorOpenTime);
            state = doorOpen;
        }
        
        break;
        
    default:
        printf("  \7WAT at %s:%d: Releasing stop button in state \"%s\" makes no sense!\n",
            __FUNCTION__, __LINE__,
            State_toString(state)
        );
    }
    
    printf("  New state: %s\n", State_toString(state));
}




void fsm_event_arrivedAtFloor(int newFloor){
    printf("%s(%d)\n  State:     %-15.15s |  Floor: %d  |  Dirn: %d\n",  
        __FUNCTION__, newFloor,
        State_toString(state), floor, dirn
    );
    
    floor = newFloor;
    
    elev_set_floor_indicator(floor);
    
    switch(state){
    case init:
        dirn = DIRN_STOP;
        elev_set_motor_direction(dirn);
        state = idle;
        break;
    case moving:
        if(orders_shouldStop(floor, dirn)){
            elev_set_motor_direction(DIRN_STOP);
            elev_set_door_open_lamp(1);
            orders_clearOrdersAt(floor);
            timer_start(doorOpenTime);
            setAllLights();
            state = doorOpen;
        } else {
            departDirn = dirn;
        }
        break;
    default:
        printf("  \7WAT at %s:%d: Arriving at a floor in state \"%s\" makes no sense!\n",
            __FUNCTION__, __LINE__,
            State_toString(state)
        );
    }
    
    printf("  New state: %s\n", State_toString(state));  
}




void fsm_event_timerHasTimedOut(void){
    printf("%s\n  State:     %-15.15s |  Floor: %d  |  Dirn: %d\n",  
        __FUNCTION__, 
        State_toString(state), floor, dirn
    );
    
    switch(state){
    case doorOpen:
        dirn = orders_chooseDirn(floor, dirn);
        
        elev_set_door_open_lamp(0);
        elev_set_motor_direction(dirn);
        
        if(dirn == DIRN_STOP){
            state = idle;
        } else {
            state = moving;
            departDirn = dirn;
        }
        
        break;
                
    default:
        printf("  \7WAT at %s:%d: Timer timing out in state \"%s\" makes no sense!\n",
            __FUNCTION__, __LINE__,
            State_toString(state)
        );
    }
    
    printf("  New state: %s\n", State_toString(state));
}
