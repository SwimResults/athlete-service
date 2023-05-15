# athlete-service

- sportlers with YoB and Club
- Teams

## Models

### Athlete

- id
- name
- year
- team (id)

### Team

- id
- name
- aliases


## API Endpoints

    [GIN-debug] GET    /athlete                  --> sr-athlete/athlete-service/controller.getAthletes (3 handlers)
    [GIN-debug] GET    /athlete/:id              --> sr-athlete/athlete-service/controller.getAthlete (3 handlers)
    [GIN-debug] GET    /athlete/meet/:meet_id    --> sr-athlete/athlete-service/controller.getAthleteByMeeting (3 handlers)
    [GIN-debug] DELETE /athlete/:id              --> sr-athlete/athlete-service/controller.removeAthlete (3 handlers)
    [GIN-debug] POST   /athlete                  --> sr-athlete/athlete-service/controller.addAthlete (3 handlers)
    [GIN-debug] POST   /athlete/participation    --> sr-athlete/athlete-service/controller.addParticipation (3 handlers)
    [GIN-debug] PUT    /athlete                  --> sr-athlete/athlete-service/controller.updateAthlete (3 handlers)
    [GIN-debug] HEAD   /athlete                  --> sr-athlete/athlete-service/controller.getAthletes (3 handlers)
    [GIN-debug] HEAD   /athlete/:id              --> sr-athlete/athlete-service/controller.getAthlete (3 handlers)
    [GIN-debug] GET    /team                     --> sr-athlete/athlete-service/controller.getTeams (3 handlers)
    [GIN-debug] GET    /team/:id                 --> sr-athlete/athlete-service/controller.getTeam (3 handlers)
    [GIN-debug] GET    /team/meet/:meet_id       --> sr-athlete/athlete-service/controller.getTeamsByMeeting (3 handlers)
    [GIN-debug] POST   /team                     --> sr-athlete/athlete-service/controller.addTeam (3 handlers)
    [GIN-debug] HEAD   /team                     --> sr-athlete/athlete-service/controller.getTeams (3 handlers)
    [GIN-debug] HEAD   /team/:id                 --> sr-athlete/athlete-service/controller.getTeam (3 handlers)
    [GIN-debug] GET    /actuator                 --> sr-athlete/athlete-service/controller.actuator (3 handlers)
