export class Room {
    constructor(data) {
        this.title = data.title;
        this.color = data.color;
        this.type = data.type;
        this.x = data.x;
        this.y = data.y;
        this.z = data.z;
        this.north = data.north;
        this.south = data.south;
        this.east = data.east;
        this.west = data.west;
        this.up = data.up;
        this.down = data.down;
    }
}