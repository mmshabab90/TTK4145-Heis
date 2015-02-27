#define _BSD_SOURCE

#include <stdio.h>
#include <unistd.h>

#include "elev.h"
#include "timer.h"
#include "fsm.h"



int main(void){

    elev_init();
    fsm_init();

    int prevButtons[N_FLOORS][3] = {{0}};
    int prevStopButton = 0;
    int prevFloor = -1;
    
    if(elev_get_floor_sensor_signal() == -1){
        elev_set_motor_direction(DIRN_DOWN);
    }
    
    while(1){    
        for(int floor = 0; floor < N_FLOORS; floor++){
            for(int btn = 0; btn < 3; btn++){
                if( (btn == BUTTON_CALL_UP      && floor == N_FLOORS-1) ||
                    (btn == BUTTON_CALL_DOWN    && floor == 0))
                {
                    continue;
                }
                
                int thisBtn = elev_get_button_signal(btn, floor);
                if(thisBtn != prevButtons[floor][btn]  &&  thisBtn){
                
                    fsm_event_buttonPressed(floor, btn);
                    
                }
                prevButtons[floor][btn] = thisBtn;
            }
        }
        
        
        int stop = elev_get_stop_signal();
        if(stop != prevStopButton){
            if(stop){
                fsm_event_stopButtonPressed();
            } else {
                fsm_event_stopButtonReleased();
            }
            prevStopButton = stop;
        }
        
        
        int floor = elev_get_floor_sensor_signal();
        if(floor != prevFloor){
            if(floor != -1){
                fsm_event_arrivedAtFloor(floor);
            }
        }
        prevFloor = floor;
        
        
        if(timer_peek()){
            timer_stop(); // goes here?
            fsm_event_timerHasTimedOut();
        }
        
        usleep(10*1000);
        
    }

    return 0;
}
