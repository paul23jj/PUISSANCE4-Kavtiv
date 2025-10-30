(function () {
    function q (id) { return ducument.getElementById(id);}

    function readMeta () {
        const el = q ( 'game-meta');
        if (!el) return { lastPlayer: 0, winner: 0 };
        const lp = parseInt(el.getAttribute('data-last-player') || '0', 10) || 0;
        const w = parseInt(el.getAttribute('data-winner') || '0', 10) || 0;
        return {lastPlayer: lp, winner: w };
    }

    function play(id) {
        const a = q(id);
        if (!a) return;
        try {
            a.currentTime = 0; 
            a.play();
        }catch (e) {

        }

    }

    function armAutoplayOnce() {
        const unlock = () => {
            ['p1-snd', 'p2-snd'].forEach((id) => {
                const a = q(id);
                if (!a) return;
                a.muted = true 
                a.play()?.finally(() => { a.pause(); a.currentTime = 0; a.muted =false; });
            });

            starter && starter.removeEventListener('click' , unlock);
            board && board.removeEventListener('click' , unlock);
            document.removeEventListener('keydown' , unlock);
        };

        const starter = q('starterBtn');
        const board = q('board');
        starter && starter.addEventListener('click', unlock, {once: true});
        board && board.addEventListener('click', unlock, {once: true });
        document.addEventListener('keydown', unlock, { once: true });
            }
            document.addEventListener('DOMContentLoaded', function (){
                armAutoplayOnce();

                const { lastPlayer , winner } = readMeta();

                if (winner === 0 ){
                    if(lastPlayer === 1 ) play('p1-snd');
                    else if (lastPlayer === 2) play('p2-snd');
                    return;
                }

                if (winner === 1) play ('p1-snd');
                else if (winner === 2) play('p2-snd');
            })
        }
        
)
