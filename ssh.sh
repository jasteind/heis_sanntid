#/bin/bash

go build main.go

scp main student@129.241.187.144:gr9/mainrun
#scp main student@129.241.187.141:gr9/mainrun

#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/main student@129.241.187.161:gruppe16/main


#gnome-terminal --title "virtual_2: 154" -x ssh student@129.241.187.154 &
#gnome-terminal --title "virtual_3: 144" -x ssh student@129.241.187.144 &
#sudo chmod 777 mainrun
