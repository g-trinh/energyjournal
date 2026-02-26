import {
  addDays,
  daysBetween,
  formatDisplayDate,
  getDaysBetween,
  getEnergyLevels,
  getEnergyLevelsRange,
  saveEnergyLevels,
  type EnergyLevels,
} from './energyLevels'

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('energyLevels service', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns parsed levels for successful GET', async () => {
    const payload: EnergyLevels = {
      date: '2026-02-21',
      physical: 7,
      mental: 5,
      emotional: 8,
      sleepQuality: 4,
      stressLevel: 2,
      physicalActivity: 'light',
      nutrition: 'good',
      socialInteractions: 'positive',
      timeOutdoors: '30min_1hr',
      notes: 'Felt productive.',
    }
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(200, payload),
    )
    const abortController = new AbortController()

    const result = await getEnergyLevels(
      '2026-02-21',
      'token',
      abortController.signal,
    )

    expect(result).toEqual(payload)
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/energy/levels?date=2026-02-21',
      expect.objectContaining({
        method: 'GET',
        signal: abortController.signal,
      }),
    )
  })

  it('returns parsed levels for successful GET range', async () => {
    const payload: EnergyLevels[] = [
      { date: '2026-02-20', physical: 6, mental: 5, emotional: 4 },
      { date: '2026-02-21', physical: 7, mental: 6, emotional: 8 },
    ]
    const fetchSpy = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(jsonResponse(200, payload))
    const abortController = new AbortController()

    const result = await getEnergyLevelsRange(
      '2026-02-20',
      '2026-02-21',
      'token',
      abortController.signal,
    )

    expect(result).toEqual(payload)
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/energy/levels/range?from=2026-02-20&to=2026-02-21',
      expect.objectContaining({
        method: 'GET',
        signal: abortController.signal,
      }),
    )
  })

  it('throws for non-200 GET range failures', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(500, { error: 'boom' }),
    )

    await expect(
      getEnergyLevelsRange('2026-02-20', '2026-02-21', 'token'),
    ).rejects.toThrow()
  })

  it('returns null for 404 GET', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(404, { error: 'not found' }),
    )

    const result = await getEnergyLevels('2026-02-21', 'token')
    expect(result).toBeNull()
  })

  it('throws for non-404 GET failures', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(500, { error: 'boom' }),
    )

    await expect(getEnergyLevels('2026-02-21', 'token')).rejects.toThrow()
  })

  it('throws on GET network failure', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('network down'))

    await expect(getEnergyLevels('2026-02-21', 'token')).rejects.toThrow(
      'network down',
    )
  })

  it('returns parsed levels for successful PUT', async () => {
    const payload: EnergyLevels = {
      date: '2026-02-21',
      physical: 6,
      mental: 5,
      emotional: 4,
      sleepQuality: 3,
      stressLevel: 2,
      physicalActivity: 'moderate',
      nutrition: 'good',
      socialInteractions: 'positive',
      timeOutdoors: '30min_1hr',
      notes: 'Great day.',
    }
    const fetchSpy = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(jsonResponse(200, payload))

    const result = await saveEnergyLevels(payload, 'token')
    expect(result).toEqual(payload)
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/energy/levels',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify(payload),
      }),
    )
  })

  it('omits empty optional context fields on PUT payload', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(200, {
        date: '2026-02-21',
        physical: 6,
        mental: 5,
        emotional: 4,
      }),
    )

    await saveEnergyLevels(
      {
        date: '2026-02-21',
        physical: 6,
        mental: 5,
        emotional: 4,
        sleepQuality: 0,
        stressLevel: 0,
        physicalActivity: '',
        nutrition: '',
        socialInteractions: '',
        timeOutdoors: '',
        notes: '',
      },
      'token',
    )

    expect(globalThis.fetch).toHaveBeenCalledWith(
      '/api/energy/levels',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({
          date: '2026-02-21',
          physical: 6,
          mental: 5,
          emotional: 4,
        }),
      }),
    )
  })

  it('throws for non-200 PUT failures', async () => {
    const payload: EnergyLevels = {
      date: '2026-02-21',
      physical: 6,
      mental: 5,
      emotional: 4,
    }
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      jsonResponse(400, { error: 'invalid' }),
    )

    await expect(saveEnergyLevels(payload, 'token')).rejects.toThrow()
  })
})

describe('energyLevels date utilities', () => {
  it('adds and subtracts days across month boundaries', () => {
    expect(addDays('2026-01-31', 1)).toBe('2026-02-01')
    expect(addDays('2026-01-01', -1)).toBe('2025-12-31')
  })

  it('returns inclusive days between range', () => {
    expect(getDaysBetween('2026-02-09', '2026-02-11')).toEqual([
      '2026-02-09',
      '2026-02-10',
      '2026-02-11',
    ])
  })

  it('returns integer day difference', () => {
    expect(daysBetween('2026-02-09', '2026-02-23')).toBe(14)
  })

  it('formats display date', () => {
    expect(formatDisplayDate('2026-02-21')).toBe('February 21, 2026')
  })
})
