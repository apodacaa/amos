from textual import events
from textual.app import App
from textual.widgets import Footer, Header, Static


class AmosApp(App):
    """A tiny Textual app for quick prototyping.

    Controls:
    - i : increment counter
    - d : decrement counter
    - q : quit
    """

    CSS = """
    Screen {
        align: center middle;
    }
    #main {
        width: 70%;
        padding: 1 2;
        border: round white;
    }
    """

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.count = 0

    def compose(self):
        yield Header()
        yield Static(self._main_text(), id="main")
        yield Footer()

    def _main_text(self) -> str:
        return (
            "Tiny Amos prototype — minimal UI\n\n"
            "Press 'i' to increment, 'd' to decrement, 'q' to quit.\n\n"
            f"Count: {self.count}"
        )

    async def on_key(self, event: events.Key) -> None:
        key = event.key.lower()
        updated = False

        if key == "i":
            self.count += 1
            updated = True
        elif key == "d":
            self.count -= 1
            updated = True
        elif key == "q":
            await self.action_quit()
            return

        if updated:
            main = self.query_one("#main", Static)
            main.update(self._main_text())
