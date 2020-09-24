namespace App {
    function rand(max: number) : number {
        return Math.floor(Math.random() * Math.floor(max));
    }

    const elements = {
        'avatar': {
            'pic': document.getElementById('avatar_status_pic'), 
            'hp': document.getElementById('avatar_status_hp'),
            'mp': document.getElementById('avatar_status_mp'),
        },
        'status': {
            'score': document.getElementById('avatar_status_score'),
            'turn': document.getElementById('avatar_status_turn'),
            'round': document.getElementById('avatar_status_round'),
        },
    };

    const urlParams = new URLSearchParams(window.location.search);

    const AVATAR_MAX_HP = 25;
    const AVATAR_MAX_MP = 20;

    function initGame() {
        elements.avatar.hp.innerHTML = `${AVATAR_MAX_HP}/${AVATAR_MAX_HP}`;
        elements.avatar.mp.innerHTML = `${AVATAR_MAX_MP}/${AVATAR_MAX_MP}`;

        let avatarID = urlParams.get('avatar') || rand(5);
        let avatarURL = `img/avatar/avatar${avatarID}.png`;
        elements.avatar.pic.innerHTML = `<img class="unit" src="${avatarURL}">`
    }

    export function main() {
        initGame();
    }
}

window.onload = function() { 
    App.main();
};
