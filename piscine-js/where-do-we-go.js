
import { places } from './where-do-we-go.data.js'

const getCoords = (coordString) => {
  const parts = coordString.match(/\d+/g).map(Number)
  const deg = parts[0]
  const min = parts[1]
  const sec = parts[2]
  const dir = coordString.includes('S') ? -1 : 1
  return dir * (deg + min / 60 + sec / 3600)
}

const createCompass = () => {
  const compass = document.createElement('div')
  compass.className = 'direction'
  compass.textContent = 'N'
  document.body.append(compass)
  return compass
}

const createLocationAnchor = (place) => {
  const anchor = document.createElement('a')
  anchor.className = 'location'
  anchor.textContent = `${place.name}\n${place.coordinates}`
  anchor.style.color = place.color
  anchor.href = `https://www.google.com/maps/place/${place.coordinates}`
  anchor.target = '_blank'
  return anchor
}

const updateLocationAnchor = (anchor, place) => {
  anchor.textContent = `${place.name}\n${place.coordinates}`
  anchor.style.color = place.color
  anchor.href = `https://www.google.com/maps/place/${place.coordinates}`
}

const createSections = (places) => {
  places.forEach((place, index) => {
    const section = document.createElement('section')
    section.id = index.toString()
    const imageLoc = place.name.toLowerCase().split(',')[0].split(' ').join('-')
    section.style.background = `URL('./where-do-we-go_images/${imageLoc}.jpg')`
    document.body.append(section)
  })
}

export const explore = () => {
  const orderedPlaces = [...places].sort((a, b) =>
        getCoords(b.coordinates) - getCoords(a.coordinates)
    )

  const compass = createCompass()
  let oldValue = 0
  window.addEventListener('scroll', () => {
    const newValue = window.pageYOffset
    compass.textContent = oldValue < newValue ? 'S' : 'N'
    oldValue = newValue
  })

  const locationAnchor = createLocationAnchor(orderedPlaces[0])
  document.body.append(locationAnchor)

  createSections(orderedPlaces)

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const val = Math.round(window.scrollY / window.innerHeight)
        updateLocationAnchor(locationAnchor, orderedPlaces[val])
      }
    })
  }, { threshold: [0.5] })

  document.querySelectorAll('section').forEach(section => observer.observe(section))
}
