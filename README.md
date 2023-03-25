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

    [GIN-debug] GET    /athlete                  --> sr-athlete/athlete-service/controller.getAthletes
    [GIN-debug] GET    /athlete/:id              --> sr-athlete/athlete-service/controller.getAthlete
    [GIN-debug] DELETE /athlete/:id              --> sr-athlete/athlete-service/controller.removeAthlete
    [GIN-debug] POST   /athlete                  --> sr-athlete/athlete-service/controller.addAthlete
    [GIN-debug] PUT    /athlete                  --> sr-athlete/athlete-service/controller.updateAthlete
    [GIN-debug] GET    /team                     --> sr-athlete/athlete-service/controller.getTeams
    [GIN-debug] GET    /team/:id                 --> sr-athlete/athlete-service/controller.getTeam
    [GIN-debug] POST   /team                     --> sr-athlete/athlete-service/controller.addTeam