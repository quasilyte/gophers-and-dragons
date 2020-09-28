namespace App {
declare class LZString {
    static compressToEncodedURIComponent(s: string): string;
    static decompressFromEncodedURIComponent(s: string): string;
}

declare function gominify(code: string): string;
declare function gofmt(code: string): string;
declare function evalGo(code: string): any;
declare function runSimulation(config: any, code: string): any;
declare function getCreepStats(name: string): any;
declare function getCardStats(name: string): any;

declare class Go {
    run(inst: WebAssembly.Instance);

    importObject: Record<string, WebAssembly.ModuleImports>;
};

export function main() {
    function loadGo(init: () => void) {
        const go = new Go();
        let stream = fetch('go.wasm', {
            cache: 'no-cache',
            headers: {
                'Cache-Control': 'no-cache',
            }
        });
        WebAssembly.instantiateStreaming(stream, go.importObject).then((result) => {
            go.run(result.instance);
            console.log('go.wasm is loaded');
            init();
        });
    }

    function copyToClipboard(text) {
        let el = document.createElement('textarea'); // Temp container
        el.value = text;
        el.setAttribute('readonly', '');
        el.style.position = 'absolute';
        el.style.left = '-9999px';
        document.body.appendChild(el);
        el.select();
        try {
            let ok = document.execCommand('copy');
            console.debug('copy to clipboard:', ok);
        } catch (e) {
            console.error('clipboard insertion failed', e);
        }
        document.body.removeChild(el);
    }

    function rand(max: number) : number {
        return Math.floor(Math.random() * Math.floor(max));
    }

    function updateElementText(el: HTMLElement, delta: number) {
        let val = parseInt(el.innerText, 10);
        let newVal = val + delta;
        if (newVal < 0) {
            newVal = 0;
        }
        el.innerText = `${newVal}`;
    }

    const elements = {
        'details': document.getElementById('hover_details'),
        'tactics': document.getElementById('tactics_editor') as HTMLTextAreaElement,
        'settings': document.getElementById('settings_editor') as HTMLTextAreaElement,
        'button': {
            'run': document.getElementById('button_run') as HTMLInputElement,
            'pause': document.getElementById('button_pause') as HTMLInputElement,
            'nextTurn': document.getElementById('button_next_turn') as HTMLInputElement,
            'format': document.getElementById('button_format') as HTMLInputElement,
            'share': document.getElementById('button_share') as HTMLInputElement,
        },
        'speed': document.getElementById('select_speed') as HTMLSelectElement,
        'tab': document.getElementById('select_tab') as HTMLSelectElement,
        'log': document.getElementById('log'),
        'avatar': {
            'pic': document.getElementById('avatar_status_pic') as HTMLImageElement, 
            'hp': document.getElementById('avatar_status_hp'),
            'mp': document.getElementById('avatar_status_mp'),
        },
        'creep': {
            'pic': document.getElementById('creep_status_pic') as HTMLImageElement,
            'name': document.getElementById('creep_status_name'),
            'hp': document.getElementById('creep_status_hp'),
        },
        'nextCreep': {
            'pic': document.getElementById('next_creep_status_pic') as HTMLImageElement,
            'name': document.getElementById('next_creep_status_name'),
            'hp': document.getElementById('next_creep_status_hp'),
        },
        'status': {
            'score': document.getElementById('status_score'),
            'turn': document.getElementById('status_turn'),
            'round': document.getElementById('status_round'),
        },
    };
    const cardElements = {
        'PowerAttack': document.getElementById('card_power_attack'),
        'Firebolt': document.getElementById('card_firebolt'),
        'Stun': document.getElementById('card_stun'),
        'Heal': document.getElementById('card_heal'),
        'Parry': document.getElementById('card_parry'),
    };

    const urlParams = new URLSearchParams(window.location.search);

    let gameSettings = {
        avatarHP: 40,
        avatarMP: 20,
        seed: null,
    }

    const NUM_ROUNDS = 10;
    const AVATAR_ID = urlParams.get('avatar') || rand(5);

    const cardDescriptions = {
        'Attack': 'Simple offensive action',
        'MagicArrow': 'Weak, but predictable way of inflicting damage',
        'Firebolt': 'Deals fire magic damage',
        'PowerAttack': 'Improved version of the Attack card',
        'Stun': 'Make the enemy skip a couple of turns',
        'Retreat': 'Avoid the current encounter',
        'Rest': 'Heal minor wounds by resting',
        'Heal': 'Heal wounds using magic',
        'Parry': 'Reflect a melee (non-ranged) back to the enemy',
    };

    let paused = false;
    let currentSimulationInterval = null;
    let currentSimulationPlayer: SimulationPlayer = null;

    function setCreep(name: string, hp: number) {
        elements.creep.pic.src = `img/creep/${name}.png`;
        elements.creep.name.innerText = name;
        elements.creep.hp.innerText = hp.toString();
    }

    function setNextCreep(name: string, hp: number) {
        elements.nextCreep.pic.src = `img/creep/${name}.png`;
        elements.nextCreep.name.innerText = name;
        elements.nextCreep.hp.innerText = hp.toString();
    }

    function encodeCodeURI(code: string): string {
        console.log("code length: %d", code.length);
        code = gominify(code);
        console.log("code minformat length: %d", code.length);
        code = LZString.compressToEncodedURIComponent(code);
        console.log("code minify+compress length: %d", code.length);
        return code;
    }

    function decodeCodeURI(uri: string): string {
        let code = LZString.decompressFromEncodedURIComponent(uri);
        return gofmt(code);
    }

    function shareURL() {
        let code = elements.tactics.value;
        let codeURI = encodeCodeURI(code);

        let site = 'https://quasilyte.dev/gophers-and-dragons/game.html';
        if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
            site = `http://${location.host}/game.html`;
        }

        if (codeURI.length > 1800) {
            return '';
        }

        var params = [];
        params.push(`avatar=${AVATAR_ID}`);
        params.push(`code=${codeURI}`);
        return site + '?' + params.join('&');
    }

    function resetPage() {
        applySettings();

        // Set spell counts to 0.
        for (const key in cardElements) {
            cardElements[key].innerText = '0';
        }
        // Set status counters to 0.
        for (const key in elements.status) {
            elements.status[key].innerText = '0';
        }
        elements.status.score.classList.remove('text-green');
        // Reset hero.
        elements.avatar.hp.innerText = `${gameSettings.avatarHP}`;
        elements.avatar.mp.innerText = `${gameSettings.avatarMP}`;
        elements.avatar.pic.src = `img/avatar/avatar${AVATAR_ID}.png`;
        // Set the initial creeps.
        setCreep('Cheepy', getCreepStats('Cheepy').maxHP);
        setNextCreep('Imp', getCreepStats('Imp').maxHP);
        // Clear the game logs.
        elements.log.innerText = '';
    }

    function renderCreepDetails(name: string, stats: any) {
        let details = `
            ${name} creep info<br>
            <br>
            HP: ${stats.maxHP}<br>
            Damage: ${stats.damage[0]}-${stats.damage[1]}<br>
            Score reward: ${stats.scoreReward}<br>
            Cards reward: ${stats.cardsReward}<br>`;
        if (stats.traits.length != 0) {
            details += `Traits: ${stats.traits.join(', ')}<br>`;
        }
        elements.details.innerHTML = details;
    }

    function renderCardDetails(name: string, stats: any) {
        let target = stats.isOffensive ? 'creep' : 'avatar';
        let details = `
            ${name} card info<br>
            <br>
            ${cardDescriptions[name]}.<br>
            <br>
            MP cost: ${stats.mp}<br>`;
        if (stats.power[1] != 0) {
            let power = (stats.power[0] == stats.power[1]) ?
                `${stats.power[0]}` : 
                `${stats.power[0]}-${stats.power[1]}`;
            details += `Power: ${power} (${stats.effect})<br>`
        }
        elements.details.innerHTML = details;
    }

    function handlePause() {
        paused = !paused;
        if (paused) {
            elements.button.pause.value = 'Resume';
        } else {
            elements.button.pause.value = 'Pause';
        }
    }

    function insertTextFirefox(field: HTMLTextAreaElement | HTMLInputElement, text: string): void {
        // Found on https://www.everythingfrontend.com/posts/insert-text-into-textarea-at-cursor-position.html 🎈
        field.setRangeText(
            text,
            field.selectionStart || 0,
            field.selectionEnd || 0,
            'end' // Without this, the cursor is either at the beginning or `text` remains selected
        );

        field.dispatchEvent(new InputEvent('input', {
            data: text,
            inputType: 'insertText',
            isComposing: false // TODO: fix @types/jsdom, this shouldn't be required
        }));
    }

    // Inserts `text` at the cursor’s position, replacing any selection, with **undo** support and by firing the `input` event.
    function insertText(field: HTMLTextAreaElement | HTMLInputElement, text: string): void {
        const document = field.ownerDocument!;
        const initialFocus = document.activeElement;
        if (initialFocus !== field) {
            field.focus();
        }

        if (!document.execCommand('insertText', false, text)) {
            insertTextFirefox(field, text);
        }

        if (initialFocus === document.body) {
            field.blur();
        } else if (initialFocus instanceof HTMLElement && initialFocus !== field) {
            initialFocus.focus();
        }
    }

    class SimulationPlayer {
        nextAction: number = 0;
        actions: any[][];

        constructor(actions: any[][]) {
            this.actions = actions;
        }

        canPlayTurn(): boolean {
            return this.nextAction < this.actions.length;
        }

        playTurn() {
            for (let i = this.nextAction; i < this.actions.length; i++) {
                this.nextAction++;
                let a = this.actions[i];
                if (a[0] == 'wait') {
                    updateElementText(elements.status.turn, 1);
                    break;
                }
                const [, ...tail] = a;
                handlers[a[0]].apply(null, tail);
            }
        }
    }

    function applySettings() {
        let settingsText = elements.settings.value;

        try {
            let x = JSON.parse(settingsText)
            if (x.avatarHP) {
                gameSettings.avatarHP = x.avatarHP;
            }
            if (x.avatarMP) {
                gameSettings.avatarMP = x.avatarMP;
            }
            if (x.seed) {
                gameSettings.seed = x.seed;
            }
        } catch (e) {
            console.error("bad settings: " + e)
        }
    }

    function initGame() {
        let code = urlParams.get('code');
        if (code !== null) {
            elements.tactics.value = decodeCodeURI(code);
        }

        elements.tab.options[0].selected = true;

        if (elements.settings.value === '') {
            elements.settings.value = JSON.stringify(gameSettings, undefined, 4);
        }

        resetPage();

        elements.creep.pic.onmouseenter = function(e) {
            let currentCreep = elements.creep.name.innerText;
            if (currentCreep === 'None') {
                return;
            }
            let creepStats = getCreepStats(currentCreep);
            renderCreepDetails(currentCreep, creepStats);
        };
        elements.nextCreep.pic.onmouseenter = function(e) {
            let nextCreep = elements.nextCreep.name.innerText;
            if (nextCreep === 'None') {
                return;
            }
            let creepStats = getCreepStats(nextCreep);
            renderCreepDetails(nextCreep, creepStats);
        };

        elements.tab.onchange = function(e) {
            let selected = elements.tab.options[elements.tab.selectedIndex].value;
            if (selected === 'tab_tactics') {
                elements.tactics.style.display = '';
                elements.settings.style.display = 'none';
            } else if (selected === 'tab_settings') {
                elements.tactics.style.display = 'none';
                elements.settings.style.display = '';
            }
        };

        let cardLabels = document.getElementsByClassName('card');
        let cardDetailsHandler = function() {
            let cardName = this.innerText;
            let cardStats = getCardStats(cardName);
            renderCardDetails(cardName, cardStats);
        };
        for (let i = 0; i < cardLabels.length; i++) {
            cardLabels[i].addEventListener('mouseenter', cardDetailsHandler, false);
        }

        elements.button.nextTurn.onclick = function(e) {
            if (paused && currentSimulationPlayer.canPlayTurn()) {
                currentSimulationPlayer.playTurn();
            }
        };

        elements.button.format.onclick = function(e) {
            let code = elements.tactics.value;
            let result = gofmt(code);
            if (result.startsWith('error:')) {
                console.error("gofmt: %s", result);
            } else {
                elements.tactics.value = result;
            }
        };

        elements.button.run.onclick = function(e) {
            if (currentSimulationInterval) {
                clearInterval(currentSimulationInterval);
                currentSimulationInterval = null;
            }
            resetPage();
            let config = {}
            config["avatarHP"] = gameSettings.avatarHP;
            config["avatarMP"] = gameSettings.avatarMP;
            config["rounds"] = NUM_ROUNDS;
            config["seed"] = gameSettings.seed;
            let code = elements.tactics.value;
            let actions = runSimulation(config, code);
            let speed = parseInt(elements.speed.options[elements.speed.selectedIndex].value, 10);
            currentSimulationPlayer = new SimulationPlayer(actions);
            console.log('starting applyActions with speed=%d', speed);
            console.log('actions:', actions);
            applyActions(speed, currentSimulationPlayer);
        };

        document.addEventListener('keyup', function(e) {
            let textareaFocused = (elements.tactics === document.activeElement);
            if (e.code === 'Space' && !textareaFocused) {
                e.preventDefault();
                handlePause();
            }
        });
        document.addEventListener('keydown', function(e) {
            let textareaFocused = (elements.tactics === document.activeElement);
            if (e.code === 'Tab' && textareaFocused) {
                e.preventDefault();
                insertText(elements.tactics, '    ');
            }
        });

        elements.button.share.onclick = function (e) {
            let url = shareURL();
            if (url) {
                copyToClipboard(url);
            } else {
                alert('Your code is too big to be shared');
            }
        };

        elements.button.pause.onclick = function(e) {
            handlePause();
        };
    }

    const handlers = {
        victory: function() {
            elements.status.score.classList.add('text-green');
        },
        defeat: function() {
            elements.avatar.pic.src = `img/dead_avatar/avatar${AVATAR_ID}.png`;
        },
        log: function(message: string) {
            elements.log.innerHTML += `${message}<br>`;
            elements.log.scrollTop = elements.log.scrollHeight;
        },
        redLog: function(message: string) {
            elements.log.innerHTML += `<span class="text-red">${message}</span><br>`;
            elements.log.scrollTop = elements.log.scrollHeight;
        },
        greenLog: function(message: string) {
            elements.log.innerHTML += `<span class="text-green">${message}</span><br>`;
            elements.log.scrollTop = elements.log.scrollHeight;
        },
        changeCardCount: function(name: string, delta: number) {
            updateElementText(cardElements[name], delta);
        },
        nextRound: function() {
            if (parseInt(elements.status.round.innerText, 10) != NUM_ROUNDS) {
                updateElementText(elements.status.round, 1);
            }
        },
        updateScore: function(delta: number) {
            updateElementText(elements.status.score, delta);
        },
        updateHP: function(delta: number) {
            updateElementText(elements.avatar.hp, delta);
        },
        updateMP: function(delta: number) {
            updateElementText(elements.avatar.mp, delta);
        },
        updateCreepHP: function(delta: number) {
            updateElementText(elements.creep.hp, delta);
        },
        setCreep: function(name: string, hp: number) {
            setCreep(name, hp);
        },
        setNextCreep: function(name: string, hp: number) {
            setNextCreep(name, hp);
        },
    };

    function applyActions(interval: number, player: SimulationPlayer) {
        currentSimulationInterval = setInterval(function() {
            if (paused) {
                return;
            }
            if (player.canPlayTurn()) {
                player.playTurn();
            } else {
                clearInterval(currentSimulationInterval);
                console.log('applied %d actions', player.actions.length);
                currentSimulationInterval = null;
            }
        }, interval);
    }

    loadGo(initGame);
}

} // namespace App

window.onload = function() { 
    App.main();
};
