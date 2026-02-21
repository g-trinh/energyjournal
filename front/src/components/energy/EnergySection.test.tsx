import { render } from '@testing-library/react'
import EnergySection from './EnergySection'

describe('EnergySection', () => {
  it('matches snapshot at value 5 with divider', () => {
    const { asFragment } = render(
      <EnergySection
        label="Physical"
        color="#c4826d"
        value={5}
        onChange={vi.fn()}
        showDivider
      />,
    )

    expect(asFragment()).toMatchSnapshot()
  })

  it('matches snapshot at value 0 without divider', () => {
    const { asFragment } = render(
      <EnergySection
        label="Emotional"
        color="#8fa58b"
        value={0}
        onChange={vi.fn()}
        showDivider={false}
      />,
    )

    expect(asFragment()).toMatchSnapshot()
  })
})
