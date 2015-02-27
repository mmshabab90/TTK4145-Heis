
#include <stdio.h>

#include "orders.h"


int orders[N_FLOORS][3];


static int ordersAbove(int floor){
    for(int f = floor+1; f < N_FLOORS; f++){
        for(int btn = 0; btn < 3; btn++){
            if(orders[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}


static int ordersBelow(int floor){
    for(int f = 0; f < floor; f++){
        for(int btn = 0; btn < 3; btn++){
            if(orders[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}


static int hasOrders(void){
    for(int f = 0; f < N_FLOORS; f++){
        for(int btn = 0; btn < 3; btn++){
            if(orders[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}



void orders_add(int floor, elev_button_type_t btn){
    orders[floor][btn] = 1;
}


elev_motor_direction_t orders_chooseDirn(int currFloor, elev_motor_direction_t currDirn){
    if(!hasOrders()){
        return DIRN_STOP;
    }
    switch(currDirn){
    case DIRN_DOWN:
        if(ordersBelow(currFloor)  &&  currFloor != 0){
            return DIRN_DOWN;
        } else {
            return DIRN_UP;
        }
        break;
    case DIRN_UP:
        if(ordersAbove(currFloor)  &&  currFloor != N_FLOORS-1){
            return DIRN_UP;
        } else {
            return DIRN_DOWN;
        }
        break;
    case DIRN_STOP:
        if(ordersAbove(currFloor)){
            return DIRN_UP;
        } else if(ordersBelow(currFloor)) {
            return DIRN_DOWN;
        } else {
            return DIRN_STOP;
        }
        break;
    default:
        printf("  WAT at %s:%d: Called with unexpected direction %d!\n", __FUNCTION__, __LINE__, currDirn);
            return DIRN_STOP;
        break;
    }
}


int orders_shouldStop(int floor, elev_motor_direction_t dirn){
    switch(dirn){
    case DIRN_DOWN:
        return
            orders[floor][BUTTON_CALL_DOWN] ||
            orders[floor][BUTTON_COMMAND]   ||
            floor == 0                      ||
            !ordersBelow(floor);
    case DIRN_UP:
        return
            orders[floor][BUTTON_CALL_UP]   ||
            orders[floor][BUTTON_COMMAND]   ||
            floor == N_FLOORS-1             ||
            !ordersAbove(floor);
    case DIRN_STOP:
        return
            orders[floor][BUTTON_CALL_DOWN] ||
            orders[floor][BUTTON_CALL_UP]   ||
            orders[floor][BUTTON_COMMAND];
    default:
        printf("  \7WAT at %s:%d: Called with unexpected direction %d!\n", 
            __FUNCTION__, __LINE__, dirn
        );
        return 0;    
    }
}


void orders_clearOrdersAt(int floor){
    for(int btn = 0; btn < 3; btn++){
        orders[floor][btn] = 0;
    }    
}


void orders_removeAll(void){
    for(int f = 0; f < N_FLOORS; f++){
        for(int btn = 0; btn < 3; btn++){
            orders[f][btn] = 0;
        }
    }
}


int orders_read(int floor, elev_button_type_t btn){
    return orders[floor][btn];
}
