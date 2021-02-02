const React = require('react')
const { Plus, Play } = require('react-feather')
const { dialog, BrowserWindow } = require('electron').remote
const axios = require('axios')

const Home = () => {
  const chooseWorld = async () => {
    const folder = await dialog.showOpenDialog(BrowserWindow.getFocusedWindow(), { properties: ['openDirectory'] })

    const res = await axios.post('http://localhost:8756/api/worlds', {
      path: folder.filePaths[0]
    })

    console.log(res)
  }

  return (
    <>
      <div className='h-full flex w-full mx-auto overflow-hidden'>
        <div className='flex-1 px-20 overflow-y-scroll'>
          <div className='mt-40'>
            <img src='/img/logo.svg' alt='nanomap' className='h-6 block' />
            <p className='mt-2 text-sm text-gray-600'>0.1.0-beta</p>
            <div className='mt-8'>
              <p>A simple map viewer for Minecraft: Bedrock Edition. <br /><a className='text-green-600 text-sm hover:underline' href="https://github.com/hrqsn/nanomap/">GitHub</a>, <a className='text-green-600 text-sm hover:underline' href="https://twitter.com/hrqsn">開発者</a></p>
            </div>
          </div>
        </div>
        <div className='flex-0 w-86 bg-white px-4'>
          <div className='h-full flex flex-col items-center justify-center'>
            <div className='w-full'>
              <button onClick={() => chooseWorld()} className='w-full flex space-x-2 hover:bg-gray-50 px-2 py-3 rounded focus:outline-none'>
                <Plus strokeWidth={1.25} />
                <div className='text-left'>
                  <span>ワールドを選択</span>
                  <p className='mt-2 text-sm text-gray-700'>ワールドデータを選択してマップを生成。</p>
                </div>
              </button>
            </div>
            <div className='mt-4 w-full'>
              <button onClick={() => chooseWorld()} className='w-full flex space-x-2 hover:bg-gray-50 px-2 py-3 rounded focus:outline-none'>
                <div className='w-6 h-6 flex items-center justify-center'>
                  <Play size={18} strokeWidth={1.25} />
                </div>
                <div className='text-left'>
                  <span>デモワールドで試す</span>
                  <p className='mt-2 text-sm text-gray-700'>デモワールドからマップを生成。</p>
                </div>
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default Home
