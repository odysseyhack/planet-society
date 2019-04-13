//
//  TransactionFlowViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionFlowViewController: UIPageViewController {

    // MARK: - Private properties

    private let steps = TransactionFlowStep.allCases
    private var currentStepIndex = 0
    private let transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(
            transitionStyle: .scroll,
            navigationOrientation: .horizontal,
            options: [:])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        // configure navigation bar
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown
        navigationItem.backBarButtonItem = UIBarButtonItem(title: "", style: .plain, target: nil, action: nil)

        // set initial viewcontroller
        let firstViewControllerOrNil = self.steps.map { $0.viewController }.first
        if let viewController = firstViewControllerOrNil {

            setViewControllers(
                [viewController],
                direction: .forward,
                animated: false)
        }
    }
}
