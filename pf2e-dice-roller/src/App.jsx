import React, { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [diceInput, setDiceInput] = useState("1d20");
  const [result, setResult] = useState(null);
  const [diceElements, setDiceElements] = useState([]);
  const diceSound = new Audio("/dice_roll.wav"); // Load sound file

  useEffect(() => {
    const handleKeyDown = (event) => {
      if (event.key === "Enter") {
        rollDice();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [diceInput]);

  const rollDice = (modifier = 0) => {
    diceSound.currentTime = 0; // Restart sound if already playing
    diceSound.play(); // Play dice rolling sound

    const dicePattern = /(\d*)d(\d+)([+\-]\d+)?/g;
    let total = 0;
    let match;
    const newDiceElements = [];

    while ((match = dicePattern.exec(diceInput)) !== null) {
      const numDice = match[1] ? parseInt(match[1]) : 1;
      const diceSides = parseInt(match[2]);
      const flatModifier = match[3] ? parseInt(match[3]) : 0;

      for (let i = 0; i < numDice; i++) {
        const roll = Math.floor(Math.random() * diceSides) + 1;
        total += roll;
        newDiceElements.push({
          roll,
          sides: diceSides,
          id: Math.random(),
          x: Math.random() * 80 + 10, // Random X position (10% to 90%)
          y: Math.random() * 80 + 10, // Random Y position (10% to 90%)
          rotation: Math.random() * 360, // Random rotation
        });
      }

      total += flatModifier + modifier;
    }

    setResult(total);
    setDiceElements(newDiceElements);

    // Remove dice after 3 seconds
    setTimeout(() => setDiceElements([]), 3000);
  };

  return (
    <div className="app">
      <h1 className="title">PF2e Dice Roller</h1>
      <input
        type="text"
        value={diceInput}
        onChange={(e) => setDiceInput(e.target.value)}
        placeholder="Enter dice (e.g., 2d6+3)"
      />
      <div className="buttons">
        <button className="fortune" onClick={() => rollDice(1)}>Fortune</button>
        <button className="misfortune" onClick={() => rollDice(-1)}>Misfortune</button>
      </div>
      {result !== null && <h2 className="result">Result: {result}</h2>}

      <div className="dice-container">
        {diceElements.map((dice) => (
          <div
            key={dice.id}
            className={`dice dice-${dice.sides}`}
            style={{
              left: `${dice.x}%`,
              top: `${dice.y}%`,
              transform: `rotate(${dice.rotation}deg)`,
            }}
          >
            {dice.roll}
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
