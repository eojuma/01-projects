
export const pick = () => {
  const hslDiv = document.createElement('div')
  const hueDiv = document.createElement('div')
  const lumDiv = document.createElement('div')
  hslDiv.className = 'text hsl'
  hueDiv.className = 'text hue'
  lumDiv.className = 'text luminosity'
  document.body.append(hslDiv, hueDiv, lumDiv)

  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
  svg.setAttribute('height', window.innerHeight)
  svg.setAttribute('width', window.innerWidth)

  const axisX = document.createElementNS('http://www.w3.org/2000/svg', 'line')
  const axisY = document.createElementNS('http://www.w3.org/2000/svg', 'line')
  axisX.setAttribute('id', 'axisX') 
  axisY.setAttribute('id', 'axisY') 
  axisX.style.stroke = axisY.style.stroke = 'white'
  axisX.style.strokeWidth = axisY.style.strokeWidth = '1'
  svg.append(axisX, axisY)
  document.body.append(svg)

  document.addEventListener('mousemove', (event) => {
    const x = event.clientX
    const y = event.clientY
    const hue = Math.round((x / window.innerWidth) * 360)
    const lum = Math.round((y / window.innerHeight) * 100)

    hueDiv.textContent = `hue ${hue}`
    lumDiv.textContent = `luminosity ${lum}`
    hslDiv.textContent = `hsl(${hue}, 50%, ${lum}%)`
    document.body.style.background = `hsl(${hue}, 50%, ${lum}%)`

    axisX.setAttribute('x1', x)
    axisX.setAttribute('x2', x)
    axisX.setAttribute('y1', 0)
    axisX.setAttribute('y2', window.innerHeight)

    axisY.setAttribute('x1', 0)
    axisY.setAttribute('x2', window.innerWidth)
    axisY.setAttribute('y1', y)
    axisY.setAttribute('y2', y)
  })

  document.addEventListener('click', () => {
    navigator.clipboard.writeText(hslDiv.textContent)
  })
}