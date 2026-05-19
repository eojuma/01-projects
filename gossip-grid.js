
import { gossips } from './gossip-grid.data.js'

export const grid = () => {
  const createGossipCard = (content) => {
    const card = document.createElement('div')
    card.className = 'gossip'
    card.textContent = content
    return card
  }
  const form = document.createElement('form')
  form.className = 'gossip'
  const textArea = document.createElement('textarea')
  const submitButton = document.createElement('button')
  submitButton.setAttribute('type', 'submit')
  submitButton.textContent = 'Share gossip!'
  form.append(textArea, submitButton)
  document.body.appendChild(form)

  submitButton.addEventListener('click', (e) => {
    e.preventDefault()
    const card = createGossipCard(textArea.value)
    document.body.insertBefore(card, form)
    textArea.value = ''
  })

  gossips.forEach(content => {
    document.body.appendChild(createGossipCard(content))
  })

  const createRange = (id, min, max, value, label) => {
    const input = document.createElement('input')
    input.type = 'range'
    input.id = id
    input.className = 'range'
    input.min = min
    input.max = max
    input.value = value
    const lbl = document.createElement('label')
    lbl.setAttribute('for', id)
    lbl.textContent = label
    lbl.style.color = 'white'
    return [lbl, input]
  }

  const ranges = document.createElement('div')
  ranges.className = 'ranges'
  const [widLabel, rangeWidth] = createRange('width', '200', '800', '50%', 'Width')
  const [fsLabel, rangeFontSize] = createRange('fontSize', '20', '40', '50%', 'Font Size')
  const [bgLabel, rangeBg] = createRange('background', '20', '75', '50', 'Background')
  ranges.append(widLabel, rangeWidth, fsLabel, rangeFontSize, bgLabel, rangeBg)
  document.body.append(ranges)

  const applyStyles = (property, value) => {
    document.querySelectorAll('.gossip').forEach(div => {
      div.style[property] = value
    })
  }

  rangeWidth.addEventListener('input', () => applyStyles('width', rangeWidth.value + 'px'))
  rangeFontSize.addEventListener('input', () => applyStyles('fontSize', rangeFontSize.value + 'px'))
  rangeBg.addEventListener('input', () => applyStyles('background', `hsl(280, 50%, ${rangeBg.value}%)`))
}